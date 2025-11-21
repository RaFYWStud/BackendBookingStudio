package repository

import (
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		Auth: ImplAuthRepository(db),
		Studio: ImplStudioRepository(db),
		Booking: ImplBookingRepository(db), 
	}
}
