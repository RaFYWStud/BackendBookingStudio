package repository

import (
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"gorm.io/gorm"
)

type bookingRepository struct {
    db *gorm.DB
}

func ImplBookingRepository(db *gorm.DB) contract.BookingRepository {
    return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *database.Booking) error {
    return r.db.Create(booking).Error
}

func (r *bookingRepository) FindByID(id int) (*database.Booking, error) {
    var booking database.Booking
    err := r.db.First(&booking, id).Error
    if err != nil {
        return nil, err
    }
    return &booking, nil
}

func (r *bookingRepository) FindByIDWithRelations(id int) (*database.Booking, error) {
    var booking database.Booking
    err := r.db.Preload("User").
        Preload("Studio").
        First(&booking, id).Error
    if err != nil {
        return nil, err
    }
    return &booking, nil
}

func (r *bookingRepository) FindAll(filter dto.BookingFilterRequest, userID *int) ([]database.Booking, int64, error) {
    var bookings []database.Booking
    var total int64

    query := r.db.Model(&database.Booking{}).
        Preload("User").
        Preload("Studio")

    // Apply filters
    if userID != nil {
        query = query.Where("user_id = ?", *userID)
    }

    if filter.StudioID > 0 {
        query = query.Where("studio_id = ?", filter.StudioID)
    }

    if filter.UserID > 0 {
        query = query.Where("user_id = ?", filter.UserID)
    }

    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }

    if filter.StartDate != "" {
        query = query.Where("booking_date >= ?", filter.StartDate)
    }

    if filter.EndDate != "" {
        query = query.Where("booking_date <= ?", filter.EndDate)
    }

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply sorting
    switch filter.SortBy {
    case "date_asc":
        query = query.Order("booking_date ASC, start_time ASC")
    case "date_desc":
        query = query.Order("booking_date DESC, start_time DESC")
    case "created_asc":
        query = query.Order("created_at ASC")
    case "created_desc":
        query = query.Order("created_at DESC")
    default:
        query = query.Order("created_at DESC")
    }

    // Apply pagination
    if filter.Page > 0 && filter.Limit > 0 {
        offset := (filter.Page - 1) * filter.Limit
        query = query.Offset(offset).Limit(filter.Limit)
    }

    err := query.Find(&bookings).Error
    return bookings, total, err
}

func (r *bookingRepository) Update(booking *database.Booking) error {
    return r.db.Save(booking).Error
}

func (r *bookingRepository) FindByUserID(userID int, filter dto.BookingFilterRequest) ([]database.Booking, int64, error) {
    return r.FindAll(filter, &userID)
}

func (r *bookingRepository) CountPendingBookings(userID int) (int64, error) {
    var count int64
    err := r.db.Model(&database.Booking{}).
        Where("user_id = ? AND status = ?", userID, "pending").
        Count(&count).Error
    return count, err
}

func (r *bookingRepository) FindExpiredBookings() ([]database.Booking, error) {
    var bookings []database.Booking
    now := time.Now()

    err := r.db.Where("status = ? AND dp_deadline < ?", "pending", now).
        Find(&bookings).Error

    return bookings, err
}