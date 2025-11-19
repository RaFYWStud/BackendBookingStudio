package dto

// ============= REQUEST DTOs =============

// CreateStudioRequest - Admin create new studio
type CreateStudioRequest struct {
    Name           string   `json:"name" binding:"required,min=3"`
    Description    string   `json:"description" binding:"required"`
    Location       string   `json:"location" binding:"required"`
    PricePerHour   int      `json:"price_per_hour" binding:"required,min=10000"`
    ImageURL       string   `json:"image_url" binding:"required,url"`
    Facilities     []string `json:"facilities" binding:"required"`
    OperatingHours string   `json:"operating_hours" binding:"required"` // Format: "09:00-22:00"
}

// UpdateStudioRequest - Admin update studio
type UpdateStudioRequest struct {
    Name           *string  `json:"name" binding:"omitempty,min=3"`
    Description    *string  `json:"description"`
    Location       *string  `json:"location"`
    PricePerHour   *int     `json:"price_per_hour" binding:"omitempty,min=10000"`
    ImageURL       *string  `json:"image_url" binding:"omitempty,url"`
    Facilities     []string `json:"facilities"`
    OperatingHours *string  `json:"operating_hours"`
    IsActive       *bool    `json:"is_active"`
}

// StudioFilterRequest - Query params for listing studios
type StudioFilterRequest struct {
    Location     string `form:"location"`
    MinPrice     int    `form:"min_price"`
    MaxPrice     int    `form:"max_price"`
    IsActive     *bool  `form:"is_active"`
    Search       string `form:"search"` // Search by name
    Page         int    `form:"page" binding:"min=1"`
    Limit        int    `form:"limit" binding:"min=1,max=100"`
    SortBy       string `form:"sort_by"` // price_asc, price_desc, name_asc, name_desc
}

// CheckAvailabilityRequest - Check studio availability
type CheckAvailabilityRequest struct {
    Date      string `json:"date" binding:"required"` // Format: "2025-11-20"
    StartTime string `json:"start_time" binding:"required"` // Format: "14:00"
    EndTime   string `json:"end_time" binding:"required"`   // Format: "17:00"
}

// ============= RESPONSE DTOs =============

// StudioResponse - Single studio detail
type StudioResponse struct {
    Success bool       `json:"success"`
    Data    StudioData `json:"data"`
}

// StudioListResponse - List of studios with pagination
type StudioListResponse struct {
    Success    bool         `json:"success"`
    Data       []StudioData `json:"data"`
    Pagination Pagination   `json:"pagination"`
}

// StudioData - Studio information
type StudioData struct {
    ID             int      `json:"id"`
    Name           string   `json:"name"`
    Description    string   `json:"description"`
    Location       string   `json:"location"`
    PricePerHour   int      `json:"price_per_hour"`
    ImageURL       string   `json:"image_url"`
    Facilities     []string `json:"facilities"`
    OperatingHours string   `json:"operating_hours"`
    IsActive       bool     `json:"is_active"`
    CreatedAt      string   `json:"created_at"`
    UpdatedAt      string   `json:"updated_at"`
}

// AvailabilityResponse - Studio availability check result
type AvailabilityResponse struct {
    Success   bool               `json:"success"`
    Available bool               `json:"available"`
    Message   string             `json:"message"`
    Data      *AvailabilityData  `json:"data,omitempty"`
}

// AvailabilityData - Available time slots
type AvailabilityData struct {
    StudioID       int             `json:"studio_id"`
    Date           string          `json:"date"`
    AvailableSlots []TimeSlot      `json:"available_slots"`
    BookedSlots    []BookedSlot    `json:"booked_slots"`
}

// TimeSlot - Available time range
type TimeSlot struct {
    StartTime string `json:"start_time"` // "09:00"
    EndTime   string `json:"end_time"`   // "12:00"
}

// BookedSlot - Booked time range
type BookedSlot struct {
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
    BookingID int    `json:"booking_id"`
}

// Pagination - Pagination metadata
type Pagination struct {
    CurrentPage  int   `json:"current_page"`
    PageSize     int   `json:"page_size"`
    TotalPages   int   `json:"total_pages"`
    TotalRecords int64 `json:"total_records"`
}

// CreateStudioResponse - Create studio response
type CreateStudioResponse struct {
    Success bool       `json:"success"`
    Message string     `json:"message"`
    Data    StudioData `json:"data"`
}

// UpdateStudioResponse - Update studio response
type UpdateStudioResponse struct {
    Success bool       `json:"success"`
    Message string     `json:"message"`
    Data    StudioData `json:"data"`
}

// DeleteStudioResponse - Delete studio response
type DeleteStudioResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}