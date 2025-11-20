package contract

import "github.com/RaFYWStud/BackendBookingStudio/dto"

type Service struct {
    Auth          AuthService
    Studio        StudioService
    Booking       BookingService
    PaymentMethod PaymentMethodService 
    Payment       PaymentService       
}

type AuthService interface {
    Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
    Login(req dto.LoginRequest) (*dto.LoginResponse, error)
    GetProfile(userID int) (*dto.ProfileResponse, error)
}

type StudioService interface {
    GetAllStudios(filter dto.StudioFilterRequest) (*dto.StudioListResponse, error)
    GetStudioByID(studioID int) (*dto.StudioResponse, error)
    CheckAvailability(studioID int, req dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error)
    CreateStudio(req dto.CreateStudioRequest) (*dto.CreateStudioResponse, error)
    UpdateStudio(studioID int, req dto.UpdateStudioRequest) (*dto.UpdateStudioResponse, error)
    DeleteStudio(studioID int) (*dto.DeleteStudioResponse, error)
}

type BookingService interface {
    CreateBooking(userID int, req dto.CreateBookingRequest) (*dto.CreateBookingResponse, error)
    GetMyBookings(userID int, filter dto.BookingFilterRequest) (*dto.BookingListResponse, error)
    GetBookingDetail(bookingID int, userID int, isAdmin bool) (*dto.BookingResponse, error)
    CancelBooking(bookingID int, userID int, req dto.CancelBookingRequest) (*dto.CancelBookingResponse, error)
    GetAllBookings(filter dto.BookingFilterRequest) (*dto.BookingListResponse, error)
    UpdateBookingStatus(bookingID int, req dto.UpdateBookingStatusRequest) (*dto.UpdateBookingStatusResponse, error)
}

type PaymentMethodService interface {
    GetActivePaymentMethods() (*dto.PaymentMethodListResponse, error)

    GetAllPaymentMethods() (*dto.PaymentMethodListResponse, error)
    GetPaymentMethodByID(id int) (*dto.PaymentMethodResponse, error)
    CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
    UpdatePaymentMethod(id int, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error)
    DeletePaymentMethod(id int) (*dto.DeletePaymentMethodResponse, error)
}

type PaymentService interface {
    UploadPaymentProof(userID int, bookingID int, req dto.UploadPaymentProofRequest) (*dto.UploadPaymentProofResponse, error)
    GetMyPayments(userID int, filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error)
    GetPaymentDetail(paymentID int, userID int, isAdmin bool) (*dto.PaymentResponse, error)
    
    GetAllPayments(filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error)
    GetPendingPayments(filter dto.PaymentFilterRequest) (*dto.PaymentListResponse, error)
    VerifyPayment(paymentID int, adminID int, req dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error)
}