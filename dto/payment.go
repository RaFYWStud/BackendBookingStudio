package dto

// ============= REQUEST DTOs =============

// UploadPaymentProofRequest - Customer upload payment proof
type UploadPaymentProofRequest struct {
    PaymentType     string `json:"payment_type" binding:"required"`      // "dp" or "full_payment"
    PaymentMethodID int    `json:"payment_method_id" binding:"required"`
    Amount          int    `json:"amount" binding:"required,min=1"`
    ProofURL        string `json:"proof_url" binding:"required,url"` // URL gambar bukti transfer
}

// VerifyPaymentRequest - Admin verify/reject payment
type VerifyPaymentRequest struct {
    Status string `json:"status" binding:"required"` // "verified" or "rejected"
    Reason string `json:"reason"`                    // Wajib jika rejected
}

// PaymentFilterRequest - Filter for payment list
type PaymentFilterRequest struct {
    BookingID int    `form:"booking_id"`
    Status    string `form:"status"` // pending, verified, rejected
    Type      string `form:"type"`   // dp, full_payment
    Page      int    `form:"page" binding:"min=1"`
    Limit     int    `form:"limit" binding:"min=1,max=100"`
}

// ============= RESPONSE DTOs =============

// PaymentResponse - Single payment detail
type PaymentResponse struct {
    Success bool              `json:"success"`
    Data    PaymentDetailData `json:"data"`
}

// PaymentListResponse - List of payments
type PaymentListResponse struct {
    Success    bool              `json:"success"`
    Data       []PaymentData     `json:"data"`
    Pagination Pagination        `json:"pagination"`
}

// PaymentDetailData - Complete payment information
type PaymentDetailData struct {
    ID               int                `json:"id"`
    BookingID        int                `json:"booking_id"`
    Booking          *BookingData       `json:"booking,omitempty"`
    PaymentType      string             `json:"payment_type"`
    PaymentMethodID  int                `json:"payment_method_id"`
    PaymentMethod    *PaymentMethodData `json:"payment_method,omitempty"`
    Amount           int                `json:"amount"`
    ProofURL         string             `json:"proof_url"`
    Status           string             `json:"status"`
    VerifiedBy       *UserData          `json:"verified_by,omitempty"`
    VerifiedAt       string             `json:"verified_at,omitempty"`
    RejectionReason  string             `json:"rejection_reason,omitempty"`
    PaidAt           string             `json:"paid_at,omitempty"`
    CreatedAt        string             `json:"created_at"`
    UpdatedAt        string             `json:"updated_at"`
}

// UploadPaymentProofResponse - Upload payment proof response
type UploadPaymentProofResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Data    PaymentDetailData `json:"data"`
}

// VerifyPaymentResponse - Verify payment response
type VerifyPaymentResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Data    PaymentDetailData `json:"data"`
}