package contract

import "github.com/unsrat-it-community/back-end-e-voting-2025/dto"

type Service struct {
    Auth    AuthService
    Studio  StudioService
    Booking BookingService
}

type AuthService interface {
    Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
    Login(req dto.LoginRequest) (*dto.LoginResponse, error)
    GetProfile(userID int) (*dto.ProfileResponse, error)
}

type StudioService interface {
    // Public endpoints
    GetAllStudios(filter dto.StudioFilterRequest) (*dto.StudioListResponse, error)
    GetStudioByID(studioID int) (*dto.StudioResponse, error)
    CheckAvailability(studioID int, req dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error)

    // Admin endpoints
    CreateStudio(req dto.CreateStudioRequest) (*dto.CreateStudioResponse, error)
    UpdateStudio(studioID int, req dto.UpdateStudioRequest) (*dto.UpdateStudioResponse, error)
    DeleteStudio(studioID int) (*dto.DeleteStudioResponse, error)
}

type BookingService interface {
    // Customer endpoints
    CreateBooking(userID int, req dto.CreateBookingRequest) (*dto.CreateBookingResponse, error)
    GetMyBookings(userID int, filter dto.BookingFilterRequest) (*dto.BookingListResponse, error)
    GetBookingDetail(bookingID int, userID int, isAdmin bool) (*dto.BookingResponse, error)
    CancelBooking(bookingID int, userID int, req dto.CancelBookingRequest) (*dto.CancelBookingResponse, error)

    // Admin endpoints
    GetAllBookings(filter dto.BookingFilterRequest) (*dto.BookingListResponse, error)
    UpdateBookingStatus(bookingID int, req dto.UpdateBookingStatusRequest) (*dto.UpdateBookingStatusResponse, error)
}