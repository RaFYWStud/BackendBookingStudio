package dto

// ============= REQUEST DTOs =============

type CreateBookingRequest struct {
    StudioID      int    `json:"studio_id" binding:"required"`
    BookingDate   string `json:"booking_date" binding:"required"` // YYYY-MM-DD
    StartTime     string `json:"start_time" binding:"required"`   // HH:MM
    EndTime       string `json:"end_time" binding:"required"`     // HH:MM
    DurationHours int    `json:"duration_hours" binding:"required,min=1"`
}

type UpdateBookingStatusRequest struct {
    Status     string `json:"status" binding:"required,oneof=pending confirmed completed cancelled"`
    AdminNotes string `json:"admin_notes"` // Catatan pembayaran/perubahan status
}

type BookingFilterRequest struct {
    Status     string `form:"status"`       // pending, confirmed, completed, cancelled
    StudioID   int    `form:"studio_id"`    // Filter by studio
    UserID     int    `form:"user_id"`      // Filter by user (admin only)
    StartDate  string `form:"start_date"`   // Filter from date (YYYY-MM-DD)
    EndDate    string `form:"end_date"`     // Filter to date (YYYY-MM-DD)
    SortBy     string `form:"sort_by"`      // date_asc, date_desc, created_asc, created_desc
    Page       int    `form:"page" binding:"min=1"`
    Limit      int    `form:"limit" binding:"min=1,max=100"`
}

type CancelBookingRequest struct {
    Reason string `json:"reason" binding:"required,min=10"`
}

// ============= RESPONSE DTOs =============

type BookingResponse struct {
    Success bool        `json:"success"`
    Data    BookingData `json:"data"`
}

type BookingListResponse struct {
    Success bool           `json:"success"`
    Data    []BookingData  `json:"data"`
    Meta    PaginationMeta `json:"meta"`
}

type CancelBookingResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

type CreateBookingResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    BookingData `json:"data"`
}

type UpdateBookingStatusResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    BookingData `json:"data"`
}

// ============= DATA DTOs =============

type BookingData struct {
    ID            int         `json:"id"`
    UserID        int         `json:"user_id"`
    User          *UserData   `json:"user,omitempty"`
    StudioID      int         `json:"studio_id"`
    Studio        *StudioData `json:"studio,omitempty"`
    BookingDate   string      `json:"booking_date"`
    StartTime     string      `json:"start_time"`
    EndTime       string      `json:"end_time"`
    DurationHours int         `json:"duration_hours"`
    TotalPrice    int         `json:"total_price"`
    Status        string      `json:"status"`
    AdminNotes    string      `json:"admin_notes,omitempty"`
    CreatedAt     string      `json:"created_at"`
    UpdatedAt     string      `json:"updated_at"`
}

// PaginationMeta - Metadata untuk pagination
type PaginationMeta struct {
    CurrentPage int   `json:"current_page"`
    PerPage     int   `json:"per_page"`
    Total       int64 `json:"total"`
    TotalPages  int   `json:"total_pages"`
}
