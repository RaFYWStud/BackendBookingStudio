package repository

import (
	"time"

	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/database"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
	"gorm.io/gorm"
)

type studioRepository struct {
    db *gorm.DB
}

func ImplStudioRepository(db *gorm.DB) contract.StudioRepository {
    return &studioRepository{db: db}
}

func (r *studioRepository) Create(studio *database.Studio) error {
    return r.db.Create(studio).Error
}

func (r *studioRepository) FindByID(id int) (*database.Studio, error) {
    var studio database.Studio
    err := r.db.First(&studio, id).Error
    if err != nil {
        return nil, err
    }
    return &studio, nil
}

func (r *studioRepository) FindAll(filter dto.StudioFilterRequest) ([]database.Studio, int64, error) {
    var studios []database.Studio
    var total int64

    query := r.db.Model(&database.Studio{})

    // Apply filters
    if filter.Location != "" {
        query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
    }
    if filter.MinPrice > 0 {
        query = query.Where("price_per_hour >= ?", filter.MinPrice)
    }
    if filter.MaxPrice > 0 {
        query = query.Where("price_per_hour <= ?", filter.MaxPrice)
    }
    if filter.IsActive != nil {
        query = query.Where("is_active = ?", *filter.IsActive)
    }
    if filter.Search != "" {
        query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
    }

    // Count total before pagination
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply sorting
    switch filter.SortBy {
    case "price_asc":
        query = query.Order("price_per_hour ASC")
    case "price_desc":
        query = query.Order("price_per_hour DESC")
    case "name_asc":
        query = query.Order("name ASC")
    case "name_desc":
        query = query.Order("name DESC")
    default:
        query = query.Order("created_at DESC")
    }

    // Apply pagination
    if filter.Page > 0 && filter.Limit > 0 {
        offset := (filter.Page - 1) * filter.Limit
        query = query.Offset(offset).Limit(filter.Limit)
    }

    err := query.Find(&studios).Error
    return studios, total, err
}

func (r *studioRepository) Update(studio *database.Studio) error {
    return r.db.Save(studio).Error
}

func (r *studioRepository) Delete(id int) error {
    return r.db.Delete(&database.Studio{}, id).Error
}

func (r *studioRepository) FindBookingsByDateRange(studioID int, date time.Time) ([]database.Booking, error) {
    var bookings []database.Booking
    err := r.db.Where("studio_id = ? AND booking_date = ? AND status NOT IN (?)",
        studioID,
        date.Format("2006-01-02"),
        []string{"cancelled", "expired"},
    ).Order("start_time ASC").Find(&bookings).Error

    return bookings, err
}

func (r *studioRepository) IsStudioAvailable(studioID int, date time.Time, startTime, endTime time.Time) (bool, error) {
    var count int64
    err := r.db.Model(&database.Booking{}).Where(
        "studio_id = ? AND booking_date = ? AND status NOT IN (?) AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?))",
        studioID,
        date.Format("2006-01-02"),
        []string{"cancelled", "expired"},
        endTime, startTime,
        endTime, endTime,
        startTime, endTime,
    ).Count(&count).Error

    return count == 0, err
}