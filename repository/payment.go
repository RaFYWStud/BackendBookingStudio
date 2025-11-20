package repository

import (
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"gorm.io/gorm"
)

type paymentRepository struct {
    db *gorm.DB
}

func ImplPaymentRepository(db *gorm.DB) contract.PaymentRepository {
    return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *database.Payment) error {
    return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id int) (*database.Payment, error) {
    var payment database.Payment
    err := r.db.First(&payment, id).Error
    if err != nil {
        return nil, err
    }
    return &payment, nil
}

func (r *paymentRepository) FindByIDWithRelations(id int) (*database.Payment, error) {
    var payment database.Payment
    err := r.db.
        Preload("Booking").
        Preload("Booking.User").
        Preload("Booking.Studio").
        Preload("VerifiedByUser").
        First(&payment, id).Error
    if err != nil {
        return nil, err
    }
    return &payment, nil
}

func (r *paymentRepository) FindByBookingID(bookingID int) ([]database.Payment, error) {
    var payments []database.Payment
    err := r.db.Where("booking_id = ?", bookingID).
        Order("created_at DESC").
        Find(&payments).Error
    return payments, err
}

func (r *paymentRepository) FindAll(filter dto.PaymentFilterRequest) ([]database.Payment, int64, error) {
    var payments []database.Payment
    var total int64

    query := r.db.Model(&database.Payment{}).
        Preload("Booking").
        Preload("Booking.User").
        Preload("Booking.Studio").
        Preload("VerifiedByUser")

    // Apply filters
    if filter.BookingID > 0 {
        query = query.Where("booking_id = ?", filter.BookingID)
    }

    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }

    if filter.Type != "" {
        query = query.Where("payment_type = ?", filter.Type)
    }

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    if filter.Page > 0 && filter.Limit > 0 {
        offset := (filter.Page - 1) * filter.Limit
        query = query.Offset(offset).Limit(filter.Limit)
    }

    // Order by latest
    query = query.Order("created_at DESC")

    err := query.Find(&payments).Error
    return payments, total, err
}

func (r *paymentRepository) FindPendingPayments(filter dto.PaymentFilterRequest) ([]database.Payment, int64, error) {
    filter.Status = "pending"
    return r.FindAll(filter)
}

func (r *paymentRepository) Update(payment *database.Payment) error {
    return r.db.Save(payment).Error
}

func (r *paymentRepository) CountByBookingAndType(bookingID int, paymentType string) (int64, error) {
    var count int64
    err := r.db.Model(&database.Payment{}).
        Where("booking_id = ? AND payment_type = ? AND status = ?", bookingID, paymentType, "verified").
        Count(&count).Error
    return count, err
}