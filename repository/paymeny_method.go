package repository

import (
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"gorm.io/gorm"
)

type paymentMethodRepository struct {
    db *gorm.DB
}

func ImplPaymentMethodRepository(db *gorm.DB) contract.PaymentMethodRepository {
    return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) Create(method *database.PaymentMethod) error {
    return r.db.Create(method).Error
}

func (r *paymentMethodRepository) FindByID(id int) (*database.PaymentMethod, error) {
    var method database.PaymentMethod
    err := r.db.First(&method, id).Error
    if err != nil {
        return nil, err
    }
    return &method, nil
}

func (r *paymentMethodRepository) FindAll() ([]database.PaymentMethod, error) {
    var methods []database.PaymentMethod
    err := r.db.Order("created_at DESC").Find(&methods).Error
    return methods, err
}

func (r *paymentMethodRepository) FindActive() ([]database.PaymentMethod, error) {
    var methods []database.PaymentMethod
    err := r.db.Where("is_active = ?", true).
        Order("name ASC").
        Find(&methods).Error
    return methods, err
}

func (r *paymentMethodRepository) Update(method *database.PaymentMethod) error {
    return r.db.Save(method).Error
}

func (r *paymentMethodRepository) Delete(id int) error {
    return r.db.Delete(&database.PaymentMethod{}, id).Error
}