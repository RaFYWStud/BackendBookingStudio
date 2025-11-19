package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/middleware"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/pkg/errs"
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
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