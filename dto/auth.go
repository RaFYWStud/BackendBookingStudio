package dto

// Register Request & Response
type RegisterRequest struct {
    Name            string `json:"name" binding:"required,min=2"`
    Email           string `json:"email" binding:"required,email"`
    Password        string `json:"password" binding:"required,min=6"`
    PasswordConfirm string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

type RegisterResponse struct {
    Success bool     `json:"success"`
    Message string   `json:"message"`
    Data    UserData `json:"data"`
}

// Login Request & Response
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Success bool      `json:"success"`
    Message string    `json:"message"`
    Data    LoginData `json:"data"`
}

type LoginData struct {
    Token string   `json:"token"`
    User  UserData `json:"user"`
}

// Profile Response
type ProfileResponse struct {
    Success bool     `json:"success"`
    Data    UserData `json:"data"`
}

// Common User Data
type UserData struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// Error Response (reusable)
type ErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
}