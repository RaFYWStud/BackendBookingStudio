package contract

import (
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
)

type Service struct {
    Auth          AuthService
    Studio        StudioService
    Booking       BookingService
    Email         EmailService   
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
    PatchStudio(studioID int, req dto.PatchStudioRequest) (*dto.PatchStudioResponse, error)
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

type EmailService interface {
    SendBookingCreated(booking *database.Booking) error              
    SendBookingConfirmed(booking *database.Booking) error           
    SendBookingCancelled(booking *database.Booking, reason string) error     
}