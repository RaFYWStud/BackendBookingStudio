package dto

// ============= REQUEST DTOs =============

// CreatePaymentMethodRequest - Admin create payment method
type CreatePaymentMethodRequest struct {
    Name          string `json:"name" binding:"required,min=2"`
    BankName      string `json:"bank_name" binding:"required"`
    AccountNumber string `json:"account_number" binding:"required"`
    AccountName   string `json:"account_name" binding:"required"`
}

// UpdatePaymentMethodRequest - Admin update payment method
type UpdatePaymentMethodRequest struct {
    Name          *string `json:"name" binding:"omitempty,min=2"`
    BankName      *string `json:"bank_name"`
    AccountNumber *string `json:"account_number"`
    AccountName   *string `json:"account_name"`
    IsActive      *bool   `json:"is_active"`
}

// ============= RESPONSE DTOs =============

// PaymentMethodResponse - Single payment method
type PaymentMethodResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message,omitempty"`
    Data    PaymentMethodData `json:"data"`
}

// PaymentMethodListResponse - List of payment methods
type PaymentMethodListResponse struct {
    Success bool                `json:"success"`
    Data    []PaymentMethodData `json:"data"`
}

// PaymentMethodData - Payment method info
type PaymentMethodData struct {
    ID            int    `json:"id"`
    Name          string `json:"name"`
    BankName      string `json:"bank_name"`
    AccountNumber string `json:"account_number"`
    AccountName   string `json:"account_name"`
    IsActive      bool   `json:"is_active"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
}

// DeletePaymentMethodResponse - Delete payment method response
type DeletePaymentMethodResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}