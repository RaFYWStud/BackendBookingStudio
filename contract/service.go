package contract

import "github.com/unsrat-it-community/back-end-e-voting-2025/dto"

type Service struct {
    Auth   AuthService
    Studio StudioService
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