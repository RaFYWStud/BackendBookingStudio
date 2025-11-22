package controller

import (
	"net/http"
	"strconv"

	"github.com/RaFYWStud/BackendBookingStudio/config/middleware"
	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"github.com/gin-gonic/gin"
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
        admin.PATCH("/:id", sc.patchStudio)
        admin.DELETE("/:id", sc.deleteStudio)
    }
}

// GetAllStudios godoc
// @Summary      Ambil semua studio
// @Description  Mengambil daftar semua studio dengan filter dan pagination
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        location   query     string  false  "Filter lokasi"
// @Param        min_price  query     int     false  "Harga minimal"
// @Param        max_price  query     int     false  "Harga maksimal"
// @Param        is_active  query     bool    false  "Hanya studio aktif"
// @Param        search     query     string  false  "Cari berdasarkan nama"
// @Param        page       query     int     false  "Halaman"                default(1)
// @Param        limit      query     int     false  "Jumlah data per halaman" default(10)
// @Param        sort_by    query     string  false  "Sortir (price_asc, price_desc, name_asc, name_desc)"
// @Success      200        {object}  dto.StudioListResponse
// @Failure      400        {object}  dto.ErrorResponse  "Invalid query parameters"
// @Failure      500        {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios [get]
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

// GetStudioByID godoc
// @Summary      Ambil 1 studio berdasarkan ID
// @Description  Mengambil detail studio berdasarkan ID
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID Studio"
// @Success      200  {object}  dto.StudioResponse
// @Failure      400  {object}  dto.ErrorResponse  "Invalid studio ID"
// @Failure      404  {object}  dto.ErrorResponse  "Studio not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios/{id} [get]
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

// CheckAvailability godoc
// @Summary      Cek jadwal ketersediaan studio
// @Description  Mengecek apakah studio tersedia pada waktu tertentu
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        id       path   int                         true  "ID Studio"
// @Param        payload  body   dto.CheckAvailabilityRequest true  "Data untuk pengecekan jadwal"
// @Success      200      {object} dto.AvailabilityResponse
// @Failure      400      {object} dto.ErrorResponse  "Invalid request payload / invalid studio ID"
// @Failure      500      {object} dto.ErrorResponse  "Internal server error"
// @Router       /studios/{id}/availability [post]
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

// CreateStudio godoc
// @Summary      Buat studio baru (Admin Only)
// @Description  Menambahkan studio baru oleh admin
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      dto.CreateStudioRequest  true  "Data studio baru"
// @Success      201      {object}  dto.CreateStudioResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request payload"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  dto.ErrorResponse  "Forbidden (bukan admin)"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios [post]
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

// UpdateStudio godoc
// @Summary      Update studio lengkap (Admin Only)
// @Description  Mengupdate seluruh data studio
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                    true  "ID Studio"
// @Param        payload  body      dto.UpdateStudioRequest true  "Data update studio"
// @Success      200      {object}  dto.UpdateStudioResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid studio ID / payload"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  dto.ErrorResponse  "Forbidden (bukan admin)"
// @Failure      404      {object}  dto.ErrorResponse  "Studio not found"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios/{id} [put]
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

// PatchStudio godoc
// @Summary      Patch sebagian data studio (Admin Only)
// @Description  Mengupdate sebagian field studio
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                   true  "ID Studio"
// @Param        payload  body      dto.PatchStudioRequest true  "Data patch studio"
// @Success      200      {object}  dto.PatchStudioResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid studio ID / payload"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  dto.ErrorResponse  "Forbidden (bukan admin)"
// @Failure      404      {object}  dto.ErrorResponse  "Studio not found"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios/{id} [patch]
func (sc *StudioController) patchStudio(ctx *gin.Context) {
    idParam := ctx.Param("id")
    studioID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid studio ID"))
        return
    }

    var payload dto.PatchStudioRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := sc.service.PatchStudio(studioID, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// DeleteStudio godoc
// @Summary      Hapus studio (Admin Only)
// @Description  Menghapus studio berdasarkan ID
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID Studio"
// @Success      200  {object}  dto.DeleteStudioResponse
// @Failure      400  {object}  dto.ErrorResponse  "Invalid studio ID"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse  "Forbidden (bukan admin)"
// @Failure      404  {object}  dto.ErrorResponse  "Studio not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /studios/{id} [delete]
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