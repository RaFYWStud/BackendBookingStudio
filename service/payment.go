package service

import (
	"fmt"
	"math"
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"gorm.io/gorm"
)

type paymentService struct {
    paymentRepo       contract.PaymentRepository
    bookingRepo       contract.BookingRepository
    paymentMethodRepo contract.PaymentMethodRepository
}

func ImplPaymentService(
    paymentRepo contract.PaymentRepository,
    bookingRepo contract.BookingRepository,
    paymentMethodRepo contract.PaymentMethodRepository,
) contract.PaymentService {
    return &paymentService{
        paymentRepo:       paymentRepo,
        bookingRepo:       bookingRepo,
        paymentMethodRepo: paymentMethodRepo,
    }
}

// UploadPaymentProof - Customer upload payment proof
func (s *paymentService) UploadPaymentProof(userID int, bookingID int, req dto.UploadPaymentProofRequest) (*dto.UploadPaymentProofResponse, error) {
    // 1. Verify booking exists and belongs to user
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("booking not found")
        }
        return nil, errs.InternalServerError("failed to fetch booking")
    }

    if booking.UserID != userID {
        return nil, errs.Forbidden("you can only upload payment for your own bookings")
    }

    // 2. Check booking status
    if booking.Status == database.BookingStatusCancelled {
        return nil, errs.BadRequest("cannot upload payment for cancelled booking")
    }

    if booking.Status == database.BookingStatusCompleted {
        return nil, errs.BadRequest("booking is already completed")
    }

    // 3. Validate payment type
    if req.PaymentType != "dp" && req.PaymentType != "full_payment" {
        return nil, errs.BadRequest("payment_type must be 'dp' or 'full_payment'")
    }

    // 4. Check if payment already exists for this type
    existingCount, err := s.paymentRepo.CountByBookingAndType(bookingID, req.PaymentType)
    if err != nil {
        return nil, errs.InternalServerError("failed to check existing payments")
    }

    if existingCount > 0 {
        return nil, errs.BadRequest(fmt.Sprintf("%s payment has already been verified for this booking", req.PaymentType))
    }

    // 5. Validate payment amount
    expectedAmount := 0
    if req.PaymentType == "dp" {
        expectedAmount = booking.DPAmount
        
        // Check DP deadline
        if time.Now().After(booking.DPDeadline) {
            return nil, errs.BadRequest("DP payment deadline has passed")
        }
    } else if req.PaymentType == "full_payment" {
        // Full payment hanya bisa dilakukan setelah DP verified
        dpCount, _ := s.paymentRepo.CountByBookingAndType(bookingID, "dp")
        if dpCount == 0 {
            return nil, errs.BadRequest("DP payment must be verified first before full payment")
        }
        expectedAmount = booking.RemainingAmount
    }

    if req.Amount != expectedAmount {
        return nil, errs.BadRequest(fmt.Sprintf("payment amount must be exactly %d", expectedAmount))
    }

    // 6. Verify payment method exists and active
    paymentMethod, err := s.paymentMethodRepo.FindByID(req.PaymentMethodID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment method not found")
        }
        return nil, errs.InternalServerError("failed to verify payment method")
    }

    if !paymentMethod.IsActive {
        return nil, errs.BadRequest("selected payment method is not active")
    }

    // 7. Create payment record
    now := time.Now()
    payment := &database.Payment{
        BookingID:   bookingID,
        PaymentType: database.PaymentType(req.PaymentType),
        Amount:      req.Amount,
        ProofURL:    req.ProofURL,
        Status:      database.PaymentStatusPending,
        PaidAt:      &now,
    }

    if err := s.paymentRepo.Create(payment); err != nil {
        return nil, errs.InternalServerError("failed to upload payment proof")
    }

    // 8. Load relations for response
    payment, _ = s.paymentRepo.FindByIDWithRelations(payment.ID)

    message := fmt.Sprintf("%s payment proof uploaded successfully. Waiting for admin verification.", req.PaymentType)

    return &dto.UploadPaymentProofResponse{
        Success: true,
        Message: message,
        Data:    s.mapPaymentToDTO(payment),
    }, nil
}

// GetMyPayments - Customer get their payment history
func (s *paymentService) GetMyPayments(userID int, filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error) {
    // Set defaults
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }

    payments, _, err := s.paymentRepo.FindAll(filter)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch payments")
    }

    // Filter hanya payment milik user
    userPayments := []database.Payment{}
    for _, payment := range payments {
        if payment.Booking != nil && payment.Booking.UserID == userID {
            userPayments = append(userPayments, payment)
        }
    }

    paymentDataList := make([]dto.PaymentData, len(userPayments))
    for i, payment := range userPayments {
        paymentDataList[i] = s.mapPaymentToDTOSimple(&payment)
    }

    totalPages := int(math.Ceil(float64(len(userPayments)) / float64(filter.Limit)))

    return &dto.PaymentListResponse{
        Success: true,
        Data:    paymentDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: int64(len(userPayments)),
        },
    }, nil
}

// GetPaymentDetail - Get payment detail
func (s *paymentService) GetPaymentDetail(paymentID int, userID int, isAdmin bool) (*dto.PaymentResponse, error) {
    payment, err := s.paymentRepo.FindByIDWithRelations(paymentID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment not found")
        }
        return nil, errs.InternalServerError("failed to fetch payment")
    }

    // Check ownership if not admin
    if !isAdmin && payment.Booking != nil && payment.Booking.UserID != userID {
        return nil, errs.Forbidden("you don't have access to this payment")
    }

    return &dto.PaymentResponse{
        Success: true,
        Data:    s.mapPaymentToDTO(payment),
    }, nil
}

// GetAllPayments - Admin get all payments
func (s *paymentService) GetAllPayments(filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error) {
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100
    }

    payments, total, err := s.paymentRepo.FindAll(filter)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch payments")
    }

    paymentDataList := make([]dto.PaymentData, len(payments))
    for i, payment := range payments {
        paymentDataList[i] = s.mapPaymentToDTOSimple(&payment)
    }

    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &dto.PaymentListResponse{
        Success: true,
        Data:    paymentDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: total,
        },
    }, nil
}

// GetPendingPayments - Admin get pending payments
func (s *paymentService) GetPendingPayments(filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error) {
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    payments, total, err := s.paymentRepo.FindPendingPayments(filter)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch pending payments")
    }

    paymentDataList := make([]dto.PaymentData, len(payments))
    for i, payment := range payments {
        paymentDataList[i] = s.mapPaymentToDTOSimple(&payment)
    }

    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &dto.PaymentListResponse{
        Success: true,
        Data:    paymentDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: total,
        },
    }, nil
}

// VerifyPayment - Admin verify or reject payment
func (s *paymentService) VerifyPayment(paymentID int, adminID int, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error) {
    // 1. Validate request
    if req.Status != "verified" && req.Status != "rejected" {
        return nil, errs.BadRequest("status must be 'verified' or 'rejected'")
    }

    if req.Status == "rejected" && req.Reason == "" {
        return nil, errs.BadRequest("rejection reason is required when rejecting payment")
    }

    // 2. Get payment with relations
    payment, err := s.paymentRepo.FindByIDWithRelations(paymentID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment not found")
        }
        return nil, errs.InternalServerError("failed to fetch payment")
    }

    // 3. Check if payment is still pending
    if payment.Status != database.PaymentStatusPending {
        return nil, errs.BadRequest(fmt.Sprintf("payment is already %s", payment.Status))
    }

    // 4. Get booking
    booking, err := s.bookingRepo.FindByID(payment.BookingID)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch booking")
    }

    // 5. Update payment status
    now := time.Now()
    payment.Status = database.PaymentStatus(req.Status)
    payment.VerifiedBy = &adminID
    payment.VerifiedAt = &now

    if req.Status == "rejected" {
        payment.RejectionReason = req.Reason
    }

    if err := s.paymentRepo.Update(payment); err != nil {
        return nil, errs.InternalServerError("failed to update payment status")
    }

    // 6. Update booking status based on payment verification
    if req.Status == "verified" {
        if payment.PaymentType == "dp" {
            // DP verified → Booking confirmed
            booking.Status = database.BookingStatusConfirmed
        } else if payment.PaymentType == "full_payment" {
            // Full payment verified → Booking paid
            booking.Status = database.BookingStatusPaid
        }

        if err := s.bookingRepo.Update(booking); err != nil {
            return nil, errs.InternalServerError("failed to update booking status")
        }
    }

    // 7. Reload payment with updated relations
    payment, _ = s.paymentRepo.FindByIDWithRelations(paymentID)

    message := fmt.Sprintf("Payment %s successfully", req.Status)
    if req.Status == "verified" {
        message = fmt.Sprintf("Payment verified. Booking status updated to %s", booking.Status)
    }

    return &dto.VerifyPaymentResponse{
        Success: true,
        Message: message,
        Data:    s.mapPaymentToDTO(payment),
    }, nil
}

// Helper: Map payment to detailed DTO
func (s *paymentService) mapPaymentToDTO(payment *database.Payment) dto.PaymentDetailData {
    data := dto.PaymentDetailData{
        ID:          payment.ID,
        BookingID:   payment.BookingID,
        PaymentType: string(payment.PaymentType),
        Amount:      payment.Amount,
        ProofURL:    payment.ProofURL,
        Status:      string(payment.Status),
        CreatedAt:   payment.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:   payment.UpdatedAt.Format("2006-01-02 15:04:05"),
    }

    if payment.RejectionReason != "" {
        data.RejectionReason = payment.RejectionReason
    }

    if payment.PaidAt != nil {
        data.PaidAt = payment.PaidAt.Format("2006-01-02 15:04:05")
    }

    if payment.VerifiedAt != nil {
        data.VerifiedAt = payment.VerifiedAt.Format("2006-01-02 15:04:05")
    }

    // Add booking info (simplified)
    if payment.Booking != nil {
        data.Booking = &dto.BookingData{
            ID:          payment.Booking.ID,
            BookingDate: payment.Booking.BookingDate.Format("2006-01-02"),
            StartTime:   payment.Booking.StartTime.Format("15:04"),
            EndTime:     payment.Booking.EndTime.Format("15:04"),
            TotalPrice:  payment.Booking.TotalPrice,
            Status:      string(payment.Booking.Status),
        }

        if payment.Booking.Studio != nil {
            data.Booking.Studio = &dto.StudioData{
                ID:       payment.Booking.Studio.ID,
                Name:     payment.Booking.Studio.Name,
                Location: payment.Booking.Studio.Location,
            }
        }
    }

    // Add verified by info
    if payment.VerifiedByUser != nil {
        data.VerifiedBy = &dto.UserData{
            ID:    payment.VerifiedByUser.ID,
            Name:  payment.VerifiedByUser.Name,
            Email: payment.VerifiedByUser.Email,
            Role:  payment.VerifiedByUser.Role,
        }
    }

    return data
}

// Helper: Map payment to simple DTO (for lists)
func (s *paymentService) mapPaymentToDTOSimple(payment *database.Payment) dto.PaymentData {
    data := dto.PaymentData{
        ID:          payment.ID,
        PaymentType: string(payment.PaymentType),
        Amount:      payment.Amount,
        Status:      string(payment.Status),
    }

    if payment.ProofURL != "" {
        data.ProofURL = payment.ProofURL
    }

    if payment.PaidAt != nil {
        data.PaidAt = payment.PaidAt.Format("2006-01-02 15:04:05")
    }

    return data
}