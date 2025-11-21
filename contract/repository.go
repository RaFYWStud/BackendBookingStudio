package contract

import (
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
)

type Repository struct {
    Auth          AuthRepository
    Studio        StudioRepository
    Booking       BookingRepository   
}

type AuthRepository interface {
    CreateUser(user *database.User) error
    FindByEmail(email string) (*database.User, error)
    FindByID(id int) (*database.User, error)
}

type StudioRepository interface {
    Create(studio *database.Studio) error
    FindByID(id int) (*database.Studio, error)
    FindAll(filter dto.StudioFilterRequest) ([]database.Studio, int64, error)
    Update(studio *database.Studio) error
    Delete(id int) error
    FindBookingsByDateRange(studioID int, date time.Time) ([]database.Booking, error)
    IsStudioAvailable(studioID int, date time.Time, startTime, endTime time.Time) (bool, error)
}

type BookingRepository interface {
    Create(booking *database.Booking) error
    FindByID(id int) (*database.Booking, error)
    FindByIDWithRelations(id int) (*database.Booking, error)
    FindAll(filter dto.BookingFilterRequest, userID *int) ([]database.Booking, int64, error)
    Update(booking *database.Booking) error
    FindByUserID(userID int, filter dto.BookingFilterRequest) ([]database.Booking, int64, error)
    CountPendingBookings(userID int) (int64, error)
    FindExpiredBookings() ([]database.Booking, error)
}