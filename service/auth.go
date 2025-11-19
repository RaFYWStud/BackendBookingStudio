package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/pkg/errs"
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/database"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
    authRepo contract.AuthRepository
}

func ImplAuthService(authRepo contract.AuthRepository) contract.AuthService {
    return &authService{authRepo: authRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
    // Check if email already exists
    _, err := s.authRepo.FindByEmail(req.Email)
    if err == nil {
        // Email found, user already exists
        return nil, errs.BadRequest("email already registered")
    }
    
    // If error is NOT "record not found", it's a database error
    if err != gorm.ErrRecordNotFound {
        return nil, errs.InternalServerError("failed to check email availability")
    }

    // Email not found (ErrRecordNotFound), proceed with registration
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errs.InternalServerError("failed to process password")
    }

    user := &database.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     "customer",
    }

    if err := s.authRepo.CreateUser(user); err != nil {
        return nil, errs.InternalServerError("failed to create user account")
    }

    return &dto.RegisterResponse{
        Success: true,
        Message: "Registration successful",
        Data: dto.UserData{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        },
    }, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
    user, err := s.authRepo.FindByEmail(req.Email)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.Unauthorized("invalid email or password")
        }
        return nil, errs.InternalServerError("failed to authenticate user")
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errs.Unauthorized("invalid email or password")
    }

    // Generate JWT token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, errs.InternalServerError("failed to generate authentication token")
    }

    return &dto.LoginResponse{
        Success: true,
        Message: "Login successful",
        Data: dto.LoginData{
            Token: token,
            User: dto.UserData{
                ID:    user.ID,
                Name:  user.Name,
                Email: user.Email,
                Role:  user.Role,
            },
        },
    }, nil
}

func (s *authService) GetProfile(userID int) (*dto.ProfileResponse, error) {
    user, err := s.authRepo.FindByID(userID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("user not found")
        }
        return nil, errs.InternalServerError("failed to fetch user profile")
    }

    return &dto.ProfileResponse{
        Success: true,
        Data: dto.UserData{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        },
    }, nil
}

func (s *authService) generateToken(user *database.User) (string, error) {
    cfg := config.Get()

    // Token expiration menggunakan ACCESS_TOKEN_LIFETIME dari .env (dalam detik)
    expirationTime := time.Now().Add(time.Duration(cfg.AccessTokenLifeTime) * time.Second)

    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    user.Role,
        "exp":     expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // TODO: Use RSA private key from config
    // For now, use simple secret (should be replaced with proper RSA key)
    secretKey := []byte("temporary-secret-key-replace-with-rsa")

    signedToken, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}