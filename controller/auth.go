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

// Register godoc
// @Summary      Register user baru
// @Description  Membuat akun baru (default role: customer)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterRequest  true  "Data registrasi user"
// @Success      201      {object}  dto.RegisterResponse
// @Failure      400      {object}  dto.ErrorResponse    "Invalid request payload / validasi gagal"
// @Failure      500      {object}  dto.ErrorResponse    "Internal server error"
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      Login user
// @Description  Login dengan email & password, mengembalikan JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginRequest  true  "Data login user"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  dto.ErrorResponse "Invalid request payload / kredensial salah"
// @Failure      500      {object}  dto.ErrorResponse "Internal server error"
// @Router       /auth/login [post]
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

// GetProfile godoc
// @Summary      Ambil profil user login
// @Description  Mengambil data profil user berdasarkan JWT token
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ProfileResponse
// @Failure      401  {object}  dto.ErrorResponse "Unauthorized / token invalid"
// @Failure      500  {object}  dto.ErrorResponse "Internal server error"
// @Router       /auth/profile [get]
func (a *AuthController) getProfile(ctx *gin.Context) {
    rawID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    id, ok := rawID.(int)
    if !ok {
        HandlerError(ctx, errs.InternalServerError("invalid user id type"))
        return
    }

    response, err := a.service.GetProfile(id)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}