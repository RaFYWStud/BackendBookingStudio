package service

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"gorm.io/gorm"
)

type bookingService struct {
    bookingRepo  contract.BookingRepository
    studioRepo   contract.StudioRepository
    emailService contract.EmailService
}

func ImplBookingService(
    bookingRepo contract.BookingRepository,
    studioRepo contract.StudioRepository,
    emailService contract.EmailService,
) contract.BookingService {
    return &bookingService{
        bookingRepo:  bookingRepo,
        studioRepo:   studioRepo,
        emailService: emailService,
    }
}

// CreateBooking - Customer create new booking
func (s *bookingService) CreateBooking(userID int, req dto.CreateBookingRequest) (*dto.CreateBookingResponse, error) {
    // 1. Verify studio exists and active
    studio, err := s.studioRepo.FindByID(req.StudioID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("studio not found")
        }
        return nil, errs.InternalServerError("failed to verify studio")
    }

    if !studio.IsActive {
        return nil, errs.BadRequest("studio is currently inactive")
    }

    // 2. Parse and validate date & times
    bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
    if err != nil {
        return nil, errs.BadRequest("invalid booking date format, use YYYY-MM-DD")
    }

    if bookingDate.Before(time.Now().Truncate(24 * time.Hour)) {
        return nil, errs.BadRequest("cannot book studio in the past")
    }

    startTime, err := time.Parse("15:04", req.StartTime)
    if err != nil {
        return nil, errs.BadRequest("invalid start_time format, use HH:MM")
    }

    endTime, err := time.Parse("15:04", req.EndTime)
    if err != nil {
        return nil, errs.BadRequest("invalid end_time format, use HH:MM")
    }

    if endTime.Before(startTime) || endTime.Equal(startTime) {
        return nil, errs.BadRequest("end_time must be after start_time")
    }

    // 3. Check studio availability
    isAvailable, err := s.studioRepo.IsStudioAvailable(req.StudioID, bookingDate, startTime, endTime)
    if err != nil {
        return nil, errs.InternalServerError("failed to check availability")
    }

    if !isAvailable {
        return nil, errs.BadRequest("studio is not available for the selected time slot")
    }

    // 4. Calculate duration, price, and DP
    duration := int(endTime.Sub(startTime).Hours())
    if duration < 1 {
        return nil, errs.BadRequest("minimum booking duration is 1 hour")
    }

    totalPrice := duration * studio.PricePerHour
    dpAmount := int(float64(totalPrice) * 0.3)
    remainingAmount := totalPrice - dpAmount

    // 5. Set DP deadline (24 hours from now)
    dpDeadline := time.Now().Add(24 * time.Hour)

    // 6. Create booking
    booking := &database.Booking{
        UserID:          userID,
        StudioID:        req.StudioID,
        BookingDate:     bookingDate,
        StartTime:       startTime,
        EndTime:         endTime,
        DurationHours:   duration,
        TotalPrice:      totalPrice,
        DPAmount:        dpAmount,
        RemainingAmount: remainingAmount,
        DPDeadline:      dpDeadline,
        Status:          database.BookingStatusPending,
    }

    if err := s.bookingRepo.Create(booking); err != nil {
        return nil, errs.InternalServerError("failed to create booking")
    }

    booking.Studio = studio

    // ⬅️ SEND EMAIL: Booking Created (status: pending)
    bookingWithRelations, err := s.bookingRepo.FindByIDWithRelations(booking.ID)
    if err != nil {
        log.Printf("⚠️  Failed to load booking relations for email: %v", err)
    } else {
        go func() {
            if err := s.emailService.SendBookingCreated(bookingWithRelations); err != nil {
                log.Printf("❌ [Email] Failed to send booking created email to %s: %v",
                    bookingWithRelations.User.Email, err)
            } else {
                log.Printf("✅ [Email] Booking created email sent to %s (Booking #%d)",
                    bookingWithRelations.User.Email, bookingWithRelations.ID)
            }
        }()
    }

    return &dto.CreateBookingResponse{
        Success: true,
        Message: "Booking created successfully. Please check your email for payment instructions.",
        Data:    s.mapBookingToDTO(booking),
    }, nil
}

// GetMyBookings - Customer get their bookings
func (s *bookingService) GetMyBookings(userID int, filter dto.BookingFilterRequest) (*dto.BookingListResponse, error) {
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }

    bookings, total, err := s.bookingRepo.FindByUserID(userID, filter)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch bookings")
    }

    bookingDataList := make([]dto.BookingData, len(bookings))
    for i, booking := range bookings {
        bookingDataList[i] = s.mapBookingToDTO(&booking)
    }

    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &dto.BookingListResponse{
        Success: true,
        Data:    bookingDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: total,
        },
    }, nil
}

// GetBookingDetail - Get booking detail with full relations
func (s *bookingService) GetBookingDetail(bookingID int, userID int, isAdmin bool) (*dto.BookingResponse, error) {
    booking, err := s.bookingRepo.FindByIDWithRelations(bookingID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("booking not found")
        }
        return nil, errs.InternalServerError("failed to fetch booking details")
    }

    if !isAdmin && booking.UserID != userID {
        return nil, errs.Forbidden("you don't have access to this booking")
    }

    return &dto.BookingResponse{
        Success: true,
        Data:    s.mapBookingToDTOWithRelations(booking),
    }, nil
}

// CancelBooking - Customer cancel booking
func (s *bookingService) CancelBooking(bookingID int, userID int, req dto.CancelBookingRequest) (*dto.CancelBookingResponse, error) {
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("booking not found")
        }
        return nil, errs.InternalServerError("failed to fetch booking")
    }

    if booking.UserID != userID {
        return nil, errs.Forbidden("you can only cancel your own bookings")
    }

    if booking.Status == database.BookingStatusCancelled {
        return nil, errs.BadRequest("booking is already cancelled")
    }

    if booking.Status == database.BookingStatusCompleted {
        return nil, errs.BadRequest("cannot cancel completed booking")
    }

    // Calculate refund
    refundAmount := 0
    if booking.Status == database.BookingStatusPaid {
        refundAmount = int(float64(booking.TotalPrice) * 0.7) // 70% refund
    }

    // Update booking status
    now := time.Now()
    booking.Status = database.BookingStatusCancelled
    booking.CancelledAt = &now
    booking.CancellationReason = req.Reason

    if err := s.bookingRepo.Update(booking); err != nil {
        return nil, errs.InternalServerError("failed to cancel booking")
    }

    cancellation := &database.Cancellation{
        BookingID:    bookingID,
        Reason:       req.Reason,
        RefundAmount: refundAmount,
        RefundStatus: "pending",
        CancelledAt:  time.Now(),
    }

    // ⬅️ SEND EMAIL: Booking Cancelled by Customer
    bookingWithRelations, err := s.bookingRepo.FindByIDWithRelations(bookingID)
    if err != nil {
        log.Printf("⚠️  Failed to load booking relations for email: %v", err)
    } else {
        go func() {
            if err := s.emailService.SendBookingCancelled(bookingWithRelations, req.Reason); err != nil {
                log.Printf("❌ [Email] Failed to send cancellation email to %s: %v",
                    bookingWithRelations.User.Email, err)
            } else {
                log.Printf("✅ [Email] Cancellation email sent to %s (Booking #%d)",
                    bookingWithRelations.User.Email, bookingID)
            }
        }()
    }

    return &dto.CancelBookingResponse{
        Success: true,
        Message: "Booking cancelled successfully. Email notification sent.",
        Data: dto.CancellationData{
            BookingID:    cancellation.BookingID,
            Reason:       cancellation.Reason,
            RefundAmount: cancellation.RefundAmount,
            RefundStatus: cancellation.RefundStatus,
            CancelledAt:  cancellation.CancelledAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}

// GetAllBookings - Admin get all bookings
func (s *bookingService) GetAllBookings(filter dto.BookingFilterRequest) (*dto.BookingListResponse, error) {
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }

    bookings, total, err := s.bookingRepo.FindAll(filter, nil)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch bookings")
    }

    bookingDataList := make([]dto.BookingData, len(bookings))
    for i, booking := range bookings {
        bookingDataList[i] = s.mapBookingToDTOWithRelations(&booking)
    }

    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &dto.BookingListResponse{
        Success: true,
        Data:    bookingDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: total,
        },
    }, nil
}

// UpdateBookingStatus - Admin update booking status
func (s *bookingService) UpdateBookingStatus(bookingID int, req dto.UpdateBookingStatusRequest) (*dto.UpdateBookingStatusResponse, error) {
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("booking not found")
        }
        return nil, errs.InternalServerError("failed to fetch booking")
    }

    validStatuses := []string{"confirmed", "completed", "cancelled"}
    isValid := false
    for _, status := range validStatuses {
        if req.Status == status {
            isValid = true
            break
        }
    }

    if !isValid {
        return nil, errs.BadRequest("invalid status. Valid statuses: confirmed, completed, cancelled")
    }

    previousStatus := booking.Status
    booking.Status = database.BookingStatus(req.Status)

    if req.Status == "cancelled" {
        now := time.Now()
        booking.CancelledAt = &now
        if req.Reason != "" {
            booking.CancellationReason = req.Reason
        } else {
            booking.CancellationReason = "Cancelled by admin"
        }
    }

    if err := s.bookingRepo.Update(booking); err != nil {
        return nil, errs.InternalServerError("failed to update booking status")
    }

    bookingWithRelations, _ := s.bookingRepo.FindByIDWithRelations(bookingID)

    // ⬅️ SEND EMAIL: Based on status change
    go func() {
        if bookingWithRelations.User == nil {
            log.Printf("⚠️  Cannot send email: User data not loaded for booking #%d", bookingID)
            return
        }

        userEmail := bookingWithRelations.User.Email

        switch req.Status {
        case "confirmed":
            // Status changed: pending/paid -> confirmed
            if previousStatus != database.BookingStatusConfirmed {
                if err := s.emailService.SendBookingConfirmed(bookingWithRelations); err != nil {
                    log.Printf("❌ [Email] Failed to send booking confirmed email to %s: %v", userEmail, err)
                } else {
                    log.Printf("✅ [Email] Booking confirmed email sent to %s (Booking #%d)", userEmail, bookingID)
                }
            }

        case "cancelled":
            // Status changed: any -> cancelled (by admin)
            reason := req.Reason
            if reason == "" {
                reason = "Cancelled by admin"
            }
            if err := s.emailService.SendBookingCancelled(bookingWithRelations, reason); err != nil {
                log.Printf("❌ [Email] Failed to send admin cancellation email to %s: %v", userEmail, err)
            } else {
                log.Printf("✅ [Email] Admin cancellation email sent to %s (Booking #%d, Reason: %s)",
                    userEmail, bookingID, reason)
            }

        case "completed":
            // Status changed: confirmed/paid -> completed
            log.Printf("ℹ️  [Email] Booking #%d marked as completed. No email configured.", bookingID)
            // Optional: Bisa tambahkan SendBookingCompleted di email service

        default:
            log.Printf("⚠️  [Email] No notification configured for status: %s", req.Status)
        }
    }()

    return &dto.UpdateBookingStatusResponse{
        Success: true,
        Message: fmt.Sprintf("Booking status updated to %s successfully. Email notification sent to customer.", req.Status),
        Data:    s.mapBookingToDTOWithRelations(bookingWithRelations),
    }, nil
}

// Helper: Map booking to DTO (basic)
func (s *bookingService) mapBookingToDTO(booking *database.Booking) dto.BookingData {
    data := dto.BookingData{
        ID:              booking.ID,
        UserID:          booking.UserID,
        StudioID:        booking.StudioID,
        BookingDate:     booking.BookingDate.Format("2006-01-02"),
        StartTime:       booking.StartTime.Format("15:04"),
        EndTime:         booking.EndTime.Format("15:04"),
        DurationHours:   booking.DurationHours,
        TotalPrice:      booking.TotalPrice,
        DPAmount:        booking.DPAmount,
        RemainingAmount: booking.RemainingAmount,
        DPDeadline:      booking.DPDeadline.Format("2006-01-02 15:04:05"),
        Status:          string(booking.Status),
        CreatedAt:       booking.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:       booking.UpdatedAt.Format("2006-01-02 15:04:05"),
    }

    if booking.Studio != nil {
        data.Studio = &dto.StudioData{
            ID:             booking.Studio.ID,
            Name:           booking.Studio.Name,
            Location:       booking.Studio.Location,
            PricePerHour:   booking.Studio.PricePerHour,
            ImageURL:       booking.Studio.ImageURL,
            OperatingHours: booking.Studio.OperatingHours,
        }
    }

    return data
}

// Helper: Map booking to DTO (with full relations)
func (s *bookingService) mapBookingToDTOWithRelations(booking *database.Booking) dto.BookingData {
    data := s.mapBookingToDTO(booking)

    if booking.User != nil {
        data.User = &dto.UserData{
            ID:    booking.User.ID,
            Name:  booking.User.Name,
            Email: booking.User.Email,
            Role:  booking.User.Role,
        }
    }

    if len(booking.Payments) > 0 {
        data.Payments = make([]dto.PaymentData, len(booking.Payments))
        for i, payment := range booking.Payments {
            data.Payments[i] = dto.PaymentData{
                ID:          payment.ID,
                PaymentType: string(payment.PaymentType),
                Amount:      payment.Amount,
                Status:      string(payment.Status),
            }

            if payment.ProofURL != "" {
                data.Payments[i].ProofURL = payment.ProofURL
            }

            if payment.PaidAt != nil {
                data.Payments[i].PaidAt = payment.PaidAt.Format("2006-01-02 15:04:05")
            }
        }
    }

    if booking.Cancellation != nil {
        data.Cancellation = &dto.CancellationData{
            ID:           booking.Cancellation.ID,
            BookingID:    booking.Cancellation.BookingID,
            Reason:       booking.Cancellation.Reason,
            RefundAmount: booking.Cancellation.RefundAmount,
            RefundStatus: booking.Cancellation.RefundStatus,
            CancelledAt:  booking.Cancellation.CancelledAt.Format("2006-01-02 15:04:05"),
        }
    }

    return data
}