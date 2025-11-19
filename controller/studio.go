package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/middleware"
	"github.com/unsrat-it-community/back-end-e-voting-2025/config/pkg/errs"
	"github.com/unsrat-it-community/back-end-e-voting-2025/contract"
	"github.com/unsrat-it-community/back-end-e-voting-2025/dto"
)

type StudioController struct {
    service contract.StudioService
}

func (sc *StudioController) GetPrefix() string {
    return "/studios"
}

func (sc *StudioController) InitService(service *contract.Service) {
    sc.service = service.Studio
}

func (sc *StudioController) InitRoute(app *gin.RouterGroup) {
    // Public routes
    app.GET("", sc.getAllStudios)
    app.GET("/:id", sc.getStudioByID)
    app.POST("/:id/availability", sc.checkAvailability)

    // Admin-only routes
    admin := app.Group("")
    admin.Use(middleware.Auth(), middleware.AdminOnly())
    {
        admin.POST("", sc.createStudio)
        admin.PUT("/:id", sc.updateStudio)
        admin.DELETE("/:id", sc.deleteStudio)
    }
}

// GetAllStudios - GET /api/studios
// Query params: location, min_price, max_price, is_active, search, page, limit, sort_by
func (sc *StudioController) getAllStudios(ctx *gin.Context) {
    var filter dto.StudioFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    response, err := sc.service.GetAllStudios(filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// GetStudioByID - GET /api/studios/:id
func (sc *StudioController) getStudioByID(ctx *gin.Context) {
    idParam := ctx.Param("id")
    studioID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid studio ID"))
        return
    }

    response, err := sc.service.GetStudioByID(studioID)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// CheckAvailability - POST /api/studios/:id/availability
func (sc *StudioController) checkAvailability(ctx *gin.Context) {
    idParam := ctx.Param("id")
    studioID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid studio ID"))
        return
    }

    var payload dto.CheckAvailabilityRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := sc.service.CheckAvailability(studioID, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// CreateStudio - POST /api/studios (Admin Only)
func (sc *StudioController) createStudio(ctx *gin.Context) {
    var payload dto.CreateStudioRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := sc.service.CreateStudio(payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, response)
}

// UpdateStudio - PUT /api/studios/:id (Admin Only)
func (sc *StudioController) updateStudio(ctx *gin.Context) {
    idParam := ctx.Param("id")
    studioID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid studio ID"))
        return
    }

    var payload dto.UpdateStudioRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := sc.service.UpdateStudio(studioID, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// DeleteStudio - DELETE /api/studios/:id (Admin Only)
func (sc *StudioController) deleteStudio(ctx *gin.Context) {
    idParam := ctx.Param("id")
    studioID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid studio ID"))
        return
    }

    response, err := sc.service.DeleteStudio(studioID)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}