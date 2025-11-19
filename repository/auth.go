package repository

import (
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/database"
	"gorm.io/gorm"
)

type authRepository struct {
    db *gorm.DB
}

func ImplAuthRepository(db *gorm.DB) contract.AuthRepository {
    return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *database.User) error {
    return r.db.Create(user).Error
}

func (r *authRepository) FindByEmail(email string) (*database.User, error) {
    var user database.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *authRepository) FindByID(id int) (*database.User, error) {
    var user database.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}