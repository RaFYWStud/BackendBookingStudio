package contract

import (
	"time"

	"github.com/unsrat-it-community/back-end-e-voting-2025/database"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
)

type Repository struct {
    Auth   AuthRepository
    Studio StudioRepository
}

type AuthRepository interface {
    CreateUser(user *database.User) error
    FindByEmail(email string) (*database.User, error)
    FindByID(id int) (*database.User, error)
}

type StudioRepository interface {
    // Basic CRUD
    Create(studio *database.Studio) error
    FindByID(id int) (*database.Studio, error)
    FindAll(filter dto.StudioFilterRequest) ([]database.Studio, int64, error)
    Update(studio *database.Studio) error
    Delete(id int) error

    // Availability check
    FindBookingsByDateRange(studioID int, date time.Time) ([]database.Booking, error)
    IsStudioAvailable(studioID int, date time.Time, startTime, endTime time.Time) (bool, error)
}