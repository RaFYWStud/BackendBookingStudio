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

    // 3. Validate duration matches request
    calculatedDuration := int(endTime.Sub(startTime).Hours())
    if calculatedDuration != req.DurationHours {
        return nil, errs.BadRequest(fmt.Sprintf("duration mismatch: calculated %d hours, but requested %d hours", calculatedDuration, req.DurationHours))
    }

    if calculatedDuration < 1 {
        return nil, errs.BadRequest("minimum booking duration is 1 hour")
    }

    // 4. Check studio availability
    isAvailable, err := s.studioRepo.IsStudioAvailable(req.StudioID, bookingDate, startTime, endTime)
    if err != nil {
        return nil, errs.InternalServerError("failed to check availability")
    }

    if !isAvailable {
        return nil, errs.BadRequest("studio is not available for the selected time slot")
    }

    // 5. Calculate total price
    totalPrice := req.DurationHours * studio.PricePerHour

    // 6. Create booking - status: pending (menunggu pembayaran manual via WhatsApp)
    booking := &database.Booking{
        UserID:        userID,
        StudioID:      req.StudioID,
        BookingDate:   bookingDate,
        StartTime:     startTime,
        EndTime:       endTime,
        DurationHours: req.DurationHours,
        TotalPrice:    totalPrice,
        Status:        database.BookingStatusPending,
    }

    if err := s.bookingRepo.Create(booking); err != nil {
        return nil, errs.InternalServerError("failed to create booking")
    }

    // 7. Load booking with relations for email
    bookingWithRelations, err := s.bookingRepo.FindByIDWithRelations(booking.ID)
    if err != nil {
        log.Printf("⚠️  Failed to load booking relations for email: %v", err)
        bookingWithRelations = booking
        bookingWithRelations.Studio = studio
    }

    // 8. Send email notification
    go func() {
        if err := s.emailService.SendBookingCreated(bookingWithRelations); err != nil {
            log.Printf("❌ [Email] Failed to send booking created email: %v", err)
        } else {
            log.Printf("✅ [Email] Booking created email sent for Booking #%d", booking.ID)
        }
    }()

    message := fmt.Sprintf("Booking berhasil dibuat. Total pembayaran: Rp %s. Silakan hubungi admin via WhatsApp untuk instruksi pembayaran.", formatRupiah(totalPrice))

    return &dto.CreateBookingResponse{
        Success: true,
        Message: message,
        Data:    s.mapBookingToDTO(bookingWithRelations),
    }, nil
}

// GetMyBookings - Customer get their bookings
func (s *bookingService) GetMyBookings(userID int, filter dto.BookingFilterRequest) (*dto.BookingListResponse, error) {
    // Set default pagination
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
        Meta: dto.PaginationMeta{
            CurrentPage: filter.Page,
            PerPage:     filter.Limit,
            Total:       total,
            TotalPages:  totalPages,
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

    // Authorization check
    if !isAdmin && booking.UserID != userID {
        return nil, errs.Forbidden("you don't have access to this booking")
    }

    return &dto.BookingResponse{
        Success: true,
        Data:    s.mapBookingToDTOWithRelations(booking),
    }, nil
}

// GetAllBookings - Admin get all bookings
func (s *bookingService) GetAllBookings(filter dto.BookingFilterRequest) (*dto.BookingListResponse, error) {
    // Set default pagination
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
        Meta: dto.PaginationMeta{
            CurrentPage: filter.Page,
            PerPage:     filter.Limit,
            Total:       total,
            TotalPages:  totalPages,
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

    // Validate status
    validStatuses := []string{"pending", "confirmed", "completed", "cancelled"}
    isValid := false
    for _, status := range validStatuses {
        if req.Status == status {
            isValid = true
            break
        }
    }

    if !isValid {
        return nil, errs.BadRequest("invalid status. Valid statuses: pending, confirmed, completed, cancelled")
    }

    previousStatus := booking.Status
    booking.Status = database.BookingStatus(req.Status)

    // Update admin notes if provided
    if req.AdminNotes != "" {
        booking.AdminNotes = req.AdminNotes
    }

    if err := s.bookingRepo.Update(booking); err != nil {
        return nil, errs.InternalServerError("failed to update booking status")
    }

    // Reload with relations
    bookingWithRelations, err := s.bookingRepo.FindByIDWithRelations(bookingID)
    if err != nil {
        log.Printf("⚠️  Failed to reload booking: %v", err)
        bookingWithRelations = booking
    }

    // Send email notification based on status change
    go func() {
        if bookingWithRelations.User == nil {
            log.Printf("⚠️  Cannot send email: User data not loaded for booking #%d", bookingID)
            return
        }

        userEmail := bookingWithRelations.User.Email

        switch req.Status {
        case "confirmed":
            // Admin confirmed payment received
            if previousStatus != database.BookingStatusConfirmed {
                if err := s.emailService.SendBookingConfirmed(bookingWithRelations); err != nil {
                    log.Printf("❌ [Email] Failed to send booking confirmed email to %s: %v", userEmail, err)
                } else {
                    log.Printf("✅ [Email] Booking confirmed email sent to %s (Booking #%d)", userEmail, bookingID)
                }
            }

        case "cancelled":
            // Admin cancelled booking
            reason := req.AdminNotes
            if reason == "" {
                reason = "Booking cancelled by admin"
            }
            if err := s.emailService.SendBookingCancelled(bookingWithRelations, reason); err != nil {
                log.Printf("❌ [Email] Failed to send cancellation email to %s: %v", userEmail, err)
            } else {
                log.Printf("✅ [Email] Cancellation email sent to %s (Booking #%d)", userEmail, bookingID)
            }

        case "completed":
            log.Printf("ℹ️  Booking #%d marked as completed", bookingID)

        default:
            log.Printf("ℹ️  Status updated to %s for booking #%d", req.Status, bookingID)
        }
    }()

    var message string
    switch req.Status {
    case "confirmed":
        message = "Booking status updated to confirmed. Customer has been notified via email."
    case "completed":
        message = "Booking marked as completed."
    case "cancelled":
        message = "Booking cancelled. Customer has been notified via email."
    default:
        message = fmt.Sprintf("Booking status updated to %s.", req.Status)
    }

    return &dto.UpdateBookingStatusResponse{
        Success: true,
        Message: message,
        Data:    s.mapBookingToDTOWithRelations(bookingWithRelations),
    }, nil
}

// CancelBooking - Customer cancel their booking
func (s *bookingService) CancelBooking(bookingID int, userID int, req dto.CancelBookingRequest) (*dto.CancelBookingResponse, error) {
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("booking not found")
        }
        return nil, errs.InternalServerError("failed to fetch booking")
    }

    // Authorization check
    if booking.UserID != userID {
        return nil, errs.Forbidden("you can only cancel your own bookings")
    }

    // Validation
    if booking.Status == database.BookingStatusCancelled {
        return nil, errs.BadRequest("booking is already cancelled")
    }

    if booking.Status == database.BookingStatusCompleted {
        return nil, errs.BadRequest("cannot cancel completed booking")
    }

    // Update status
    booking.Status = database.BookingStatusCancelled
    booking.AdminNotes = fmt.Sprintf("Cancelled by customer. Reason: %s", req.Reason)

    if err := s.bookingRepo.Update(booking); err != nil {
        return nil, errs.InternalServerError("failed to cancel booking")
    }

    // Reload with relations
    bookingWithRelations, err := s.bookingRepo.FindByIDWithRelations(bookingID)
    if err != nil {
        log.Printf("⚠️  Failed to reload booking: %v", err)
    }

    // Send email
    go func() {
        if err := s.emailService.SendBookingCancelled(bookingWithRelations, req.Reason); err != nil {
            log.Printf("❌ [Email] Failed to send cancellation email: %v", err)
        } else {
            log.Printf("✅ [Email] Cancellation email sent for Booking #%d", bookingID)
        }
    }()

    return &dto.CancelBookingResponse{
        Success: true,
        Message: "Booking cancelled successfully. Admin has been notified.",
    }, nil
}

// ============= HELPER FUNCTIONS =============

// mapBookingToDTO - Basic mapping (untuk list)
func (s *bookingService) mapBookingToDTO(booking *database.Booking) dto.BookingData {
    data := dto.BookingData{
        ID:            booking.ID,
        UserID:        booking.UserID,
        StudioID:      booking.StudioID,
        BookingDate:   booking.BookingDate.Format("2006-01-02"),
        StartTime:     booking.StartTime.Format("15:04"),
        EndTime:       booking.EndTime.Format("15:04"),
        DurationHours: booking.DurationHours,
        TotalPrice:    booking.TotalPrice,
        Status:        string(booking.Status),
        AdminNotes:    booking.AdminNotes,
        CreatedAt:     booking.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:     booking.UpdatedAt.Format("2006-01-02 15:04:05"),
    }

    // Include studio if loaded
    if booking.Studio != nil {
        data.Studio = &dto.StudioData{
            ID:             booking.Studio.ID,
            Name:           booking.Studio.Name,
            Location:       booking.Studio.Location,
            PricePerHour:   booking.Studio.PricePerHour,
            ImageURL:       booking.Studio.ImageURL,
            Facilities:     booking.Studio.Facilities,
            OperatingHours: booking.Studio.OperatingHours,
        }
    }

    return data
}

// mapBookingToDTOWithRelations - Full mapping with relations (untuk detail)
func (s *bookingService) mapBookingToDTOWithRelations(booking *database.Booking) dto.BookingData {
    data := s.mapBookingToDTO(booking)

    // Include user info (for admin view)
    if booking.User != nil {
        data.User = &dto.UserData{
            ID:    booking.User.ID,
            Name:  booking.User.Name,
            Email: booking.User.Email,
            Role:  booking.User.Role,
        }
    }

    return data
}

// formatRupiah - Format number as Rupiah
func formatRupiah(amount int) string {
    if amount < 1000 {
        return fmt.Sprintf("%d", amount)
    }

    str := fmt.Sprintf("%d", amount)
    length := len(str)
    result := ""

    for i, char := range str {
        if i > 0 && (length-i)%3 == 0 {
            result += "."
        }
        result += string(char)
    }

    return result
}