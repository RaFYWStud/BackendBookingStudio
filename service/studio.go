package service

import (
	"fmt"
	"math"
	"time"

	"github.com/unsrat-it-community/back-end-e-voting-2025/config/pkg/errs"
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/database"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
	"gorm.io/gorm"
)

type studioService struct {
    studioRepo contract.StudioRepository
}

func ImplStudioService(studioRepo contract.StudioRepository) contract.StudioService {
    return &studioService{studioRepo: studioRepo}
}

// GetAllStudios - Get list of studios with filters and pagination
func (s *studioService) GetAllStudios(filter dto.StudioFilterRequest) (*dto.StudioListResponse, error) {
    // Set default pagination values
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }
    if filter.Limit > 100 {
        filter.Limit = 100 // Max limit
    }

    studios, total, err := s.studioRepo.FindAll(filter)
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch studios")
    }

    // Convert to DTO
    studioDataList := make([]dto.StudioData, len(studios))
    for i, studio := range studios {
        studioDataList[i] = dto.StudioData{
            ID:             studio.ID,
            Name:           studio.Name,
            Description:    studio.Description,
            Location:       studio.Location,
            PricePerHour:   studio.PricePerHour,
            ImageURL:       studio.ImageURL,
            Facilities:     studio.Facilities,
            OperatingHours: studio.OperatingHours,
            IsActive:       studio.IsActive,
            CreatedAt:      studio.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:      studio.UpdatedAt.Format("2006-01-02 15:04:05"),
        }
    }

    // Calculate pagination
    totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

    return &dto.StudioListResponse{
        Success: true,
        Data:    studioDataList,
        Pagination: dto.Pagination{
            CurrentPage:  filter.Page,
            PageSize:     filter.Limit,
            TotalPages:   totalPages,
            TotalRecords: total,
        },
    }, nil
}

// GetStudioByID - Get single studio detail
func (s *studioService) GetStudioByID(studioID int) (*dto.StudioResponse, error) {
    studio, err := s.studioRepo.FindByID(studioID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("studio not found")
        }
        return nil, errs.InternalServerError("failed to fetch studio details")
    }

    return &dto.StudioResponse{
        Success: true,
        Data: dto.StudioData{
            ID:             studio.ID,
            Name:           studio.Name,
            Description:    studio.Description,
            Location:       studio.Location,
            PricePerHour:   studio.PricePerHour,
            ImageURL:       studio.ImageURL,
            Facilities:     studio.Facilities,
            OperatingHours: studio.OperatingHours,
            IsActive:       studio.IsActive,
            CreatedAt:      studio.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:      studio.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}

// CheckAvailability - Check studio availability for specific date and time
func (s *studioService) CheckAvailability(studioID int, req dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error) {
    // Verify studio exists
    _, err := s.studioRepo.FindByID(studioID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("studio not found")
        }
        return nil, errs.InternalServerError("failed to check studio")
    }

    // Parse date
    date, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        return nil, errs.BadRequest("invalid date format, use YYYY-MM-DD")
    }

    // Parse times
    startTime, err := time.Parse("15:04", req.StartTime)
    if err != nil {
        return nil, errs.BadRequest("invalid start_time format, use HH:MM")
    }

    endTime, err := time.Parse("15:04", req.EndTime)
    if err != nil {
        return nil, errs.BadRequest("invalid end_time format, use HH:MM")
    }

    // Validate time range
    if endTime.Before(startTime) || endTime.Equal(startTime) {
        return nil, errs.BadRequest("end_time must be after start_time")
    }

    // Get all bookings for this studio on this date
    bookings, err := s.studioRepo.FindBookingsByDateRange(studioID, date)
    if err != nil {
        return nil, errs.InternalServerError("failed to check bookings")
    }

    // Check availability
    isAvailable, err := s.studioRepo.IsStudioAvailable(studioID, date, startTime, endTime)
    if err != nil {
        return nil, errs.InternalServerError("failed to verify availability")
    }

    // Build booked slots
    bookedSlots := make([]dto.BookedSlot, len(bookings))
    for i, booking := range bookings {
        bookedSlots[i] = dto.BookedSlot{
            StartTime: booking.StartTime.Format("15:04"),
            EndTime:   booking.EndTime.Format("15:04"),
            BookingID: booking.ID,
        }
    }

    // Build available slots (simplified - you can enhance this)
    var availableSlots []dto.TimeSlot
    if isAvailable {
        availableSlots = []dto.TimeSlot{
            {
                StartTime: req.StartTime,
                EndTime:   req.EndTime,
            },
        }
    }

    message := "Studio is available for the requested time"
    if !isAvailable {
        message = "Studio is not available for the requested time. Please check booked slots."
    }

    return &dto.AvailabilityResponse{
        Success:   true,
        Available: isAvailable,
        Message:   message,
        Data: &dto.AvailabilityData{
            StudioID:       studioID,
            Date:           req.Date,
            AvailableSlots: availableSlots,
            BookedSlots:    bookedSlots,
        },
    }, nil
}

// CreateStudio - Admin create new studio
func (s *studioService) CreateStudio(req dto.CreateStudioRequest) (*dto.CreateStudioResponse, error) {
    studio := &database.Studio{
        Name:           req.Name,
        Description:    req.Description,
        Location:       req.Location,
        PricePerHour:   req.PricePerHour,
        ImageURL:       req.ImageURL,
        Facilities:     database.StringArray(req.Facilities),
        OperatingHours: req.OperatingHours,
        IsActive:       true,
    }

    if err := s.studioRepo.Create(studio); err != nil {
        return nil, errs.InternalServerError("failed to create studio")
    }

    return &dto.CreateStudioResponse{
        Success: true,
        Message: "Studio created successfully",
        Data: dto.StudioData{
            ID:             studio.ID,
            Name:           studio.Name,
            Description:    studio.Description,
            Location:       studio.Location,
            PricePerHour:   studio.PricePerHour,
            ImageURL:       studio.ImageURL,
            Facilities:     studio.Facilities,
            OperatingHours: studio.OperatingHours,
            IsActive:       studio.IsActive,
            CreatedAt:      studio.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:      studio.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}

// UpdateStudio - Admin update studio
func (s *studioService) UpdateStudio(studioID int, req dto.UpdateStudioRequest) (*dto.UpdateStudioResponse, error) {
    // Find existing studio
    studio, err := s.studioRepo.FindByID(studioID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("studio not found")
        }
        return nil, errs.InternalServerError("failed to fetch studio")
    }

    // Update fields if provided
    if req.Name != nil {
        studio.Name = *req.Name
    }
    if req.Description != nil {
        studio.Description = *req.Description
    }
    if req.Location != nil {
        studio.Location = *req.Location
    }
    if req.PricePerHour != nil {
        studio.PricePerHour = *req.PricePerHour
    }
    if req.ImageURL != nil {
        studio.ImageURL = *req.ImageURL
    }
    if len(req.Facilities) > 0 {
    	studio.Facilities = database.StringArray(req.Facilities)
	}
    if req.OperatingHours != nil {
        studio.OperatingHours = *req.OperatingHours
    }
    if req.IsActive != nil {
        studio.IsActive = *req.IsActive
    }

    if err := s.studioRepo.Update(studio); err != nil {
        return nil, errs.InternalServerError("failed to update studio")
    }

    return &dto.UpdateStudioResponse{
        Success: true,
        Message: "Studio updated successfully",
        Data: dto.StudioData{
            ID:             studio.ID,
            Name:           studio.Name,
            Description:    studio.Description,
            Location:       studio.Location,
            PricePerHour:   studio.PricePerHour,
            ImageURL:       studio.ImageURL,
            Facilities:     studio.Facilities,
            OperatingHours: studio.OperatingHours,
            IsActive:       studio.IsActive,
            CreatedAt:      studio.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt:      studio.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}

// DeleteStudio - Admin delete/deactivate studio
func (s *studioService) DeleteStudio(studioID int) (*dto.DeleteStudioResponse, error) {
    // Check if studio exists
    _, err := s.studioRepo.FindByID(studioID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("studio not found")
        }
        return nil, errs.InternalServerError("failed to fetch studio")
    }

    // Soft delete (you can also do hard delete or just set IsActive = false)
    if err := s.studioRepo.Delete(studioID); err != nil {
        return nil, errs.InternalServerError("failed to delete studio")
    }

    return &dto.DeleteStudioResponse{
        Success: true,
        Message: fmt.Sprintf("Studio with ID %d has been deleted successfully", studioID),
    }, nil
}