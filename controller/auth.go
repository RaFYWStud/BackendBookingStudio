package controller

import (
	"net/http"

	"github.com/RaFYWStud/BackendBookingStudio/config/middleware"
	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
    service contract.AuthService
}

func (a *AuthController) GetPrefix() string {
    return "/auth"
}

func (a *AuthController) InitService(service *contract.Service) {
    a.service = service.Auth
}

func (a *AuthController) InitRoute(app *gin.RouterGroup) {
    // Public routes
    app.POST("/register", a.register)
    app.POST("/login", a.login)
    
    // Protected routes (require authentication)
    app.GET("/profile", middleware.Auth(), a.getProfile)
}

// Register - POST /api/auth/register
func (a *AuthController) register(ctx *gin.Context) {
    var payload dto.RegisterRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := a.service.Register(payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, response)
}

// Login - POST /api/auth/login
func (a *AuthController) login(ctx *gin.Context) {
    var payload dto.LoginRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := a.service.Login(payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// GetProfile - GET /api/auth/profile (Protected)
func (a *AuthController) getProfile(ctx *gin.Context) {
    // Get user_id from JWT middleware context
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    response, err := a.service.GetProfile(userID.(int))
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}