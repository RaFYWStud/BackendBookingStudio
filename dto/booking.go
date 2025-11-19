package dto

// ============= REQUEST DTOs =============

// CreateBookingRequest - Customer create new booking
type CreateBookingRequest struct {
    StudioID    int    `json:"studio_id" binding:"required"`
    BookingDate string `json:"booking_date" binding:"required"` // Format: "2025-11-25"
    StartTime   string `json:"start_time" binding:"required"`   // Format: "14:00"
    EndTime     string `json:"end_time" binding:"required"`     // Format: "17:00"
}

// BookingFilterRequest - Query params for listing bookings
type BookingFilterRequest struct {
    Status      string `form:"status"`       // pending, confirmed, paid, completed, cancelled
    StudioID    int    `form:"studio_id"`    // Filter by studio
    UserID      int    `form:"user_id"`      // Admin: filter by user
    StartDate   string `form:"start_date"`   // Filter from date
    EndDate     string `form:"end_date"`     // Filter to date
    Page        int    `form:"page" binding:"min=1"`
    Limit       int    `form:"limit" binding:"min=1,max=100"`
    SortBy      string `form:"sort_by"` // date_asc, date_desc, created_asc, created_desc
}

// CancelBookingRequest - Customer cancel booking
type CancelBookingRequest struct {
    Reason string `json:"reason" binding:"required,min=10"`
}

// UpdateBookingStatusRequest - Admin update booking status
type UpdateBookingStatusRequest struct {
    Status string `json:"status" binding:"required"` // confirmed, completed, cancelled
    Reason string `json:"reason"` // Optional, for cancellation
}

// ============= RESPONSE DTOs =============

// BookingResponse - Single booking detail
type BookingResponse struct {
    Success bool        `json:"success"`
    Data    BookingData `json:"data"`
}

// BookingListResponse - List of bookings with pagination
type BookingListResponse struct {
    Success    bool          `json:"success"`
    Data       []BookingData `json:"data"`
    Pagination Pagination    `json:"pagination"`
}

// BookingData - Booking information
type BookingData struct {
    ID              int               `json:"id"`
    UserID          int               `json:"user_id"`
    User            *UserData         `json:"user,omitempty"` // Include for admin
    StudioID        int               `json:"studio_id"`
    Studio          *StudioData       `json:"studio,omitempty"`
    BookingDate     string            `json:"booking_date"`
    StartTime       string            `json:"start_time"`
    EndTime         string            `json:"end_time"`
    DurationHours   int               `json:"duration_hours"`
    TotalPrice      int               `json:"total_price"`
    DPAmount        int               `json:"dp_amount"`
    RemainingAmount int               `json:"remaining_amount"`
    DPDeadline      string            `json:"dp_deadline"`
    Status          string            `json:"status"`
    Payments        []PaymentData     `json:"payments,omitempty"`
    Cancellation    *CancellationData `json:"cancellation,omitempty"`
    CreatedAt       string            `json:"created_at"`
    UpdatedAt       string            `json:"updated_at"`
}

// CreateBookingResponse - Create booking response
type CreateBookingResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    BookingData `json:"data"`
}

// CancelBookingResponse - Cancel booking response
type CancelBookingResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Data    CancellationData  `json:"data"`
}

// UpdateBookingStatusResponse - Update status response
type UpdateBookingStatusResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    BookingData `json:"data"`
}

// CancellationData - Cancellation information
type CancellationData struct {
    ID           int    `json:"id"`
    BookingID    int    `json:"booking_id"`
    Reason       string `json:"reason"`
    RefundAmount int    `json:"refund_amount"`
    RefundStatus string `json:"refund_status"`
    CancelledAt  string `json:"cancelled_at"`
}

// PaymentData - Payment information (simplified for booking response)
type PaymentData struct {
    ID          int    `json:"id"`
    PaymentType string `json:"payment_type"`
    Amount      int    `json:"amount"`
    Status      string `json:"status"`
    ProofURL    string `json:"proof_url,omitempty"`
    PaidAt      string `json:"paid_at,omitempty"`
}