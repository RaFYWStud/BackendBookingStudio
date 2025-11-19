package repository

import (
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		Auth: ImplAuthRepository(db),
		Studio: ImplStudioRepository(db),
	}
}
