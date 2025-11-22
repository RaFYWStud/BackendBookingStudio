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

type BookingController struct {
    service contract.BookingService
}

func (bc *BookingController) GetPrefix() string {
    return "/bookings"
}

func (bc *BookingController) InitService(service *contract.Service) {
    bc.service = service.Booking
}

func (bc *BookingController) InitRoute(app *gin.RouterGroup) {
    // Customer routes (require auth)
    customer := app.Group("")
    customer.Use(middleware.Auth())
    {
        customer.POST("", bc.createBooking)
        customer.GET("", bc.getMyBookings)
        customer.GET("/:id", bc.getBookingDetail)
        customer.POST("/:id/cancel", bc.cancelBooking)
    }

    // Admin routes
    admin := app.Group("/admin")
    admin.Use(middleware.Auth(), middleware.AdminOnly())
    {
        admin.GET("", bc.getAllBookings)
        admin.PUT("/:id/status", bc.updateBookingStatus)
    }
}

// CreateBooking godoc
// @Summary      Buat booking baru
// @Description  Customer membuat booking studio
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      dto.CreateBookingRequest  true  "Data booking"
// @Success      201      {object}  dto.CreateBookingResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /bookings [post]
func (bc *BookingController) createBooking(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    var payload dto.CreateBookingRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := bc.service.CreateBooking(userID.(int), payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, response)
}

// GetMyBookings godoc
// @Summary      Ambil riwayat booking user
// @Description  Mengambil semua booking milik user login
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status    query   string  false  "Filter status"
// @Param        page      query   int     false  "Halaman"
// @Param        limit     query   int     false  "Jumlah per halaman"
// @Success      200       {object} dto.BookingListResponse
// @Failure      400       {object} dto.ErrorResponse
// @Failure      401       {object} dto.ErrorResponse
// @Router       /bookings [get]
func (bc *BookingController) getMyBookings(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    var filter dto.BookingFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    response, err := bc.service.GetMyBookings(userID.(int), filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// GetBookingDetail godoc
// @Summary      Ambil detail 1 booking
// @Description  Customer/Admin bisa lihat detail booking
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id        path     int  true  "ID Booking"
// @Success      200       {object} dto.BookingResponse
// @Failure      400       {object} dto.ErrorResponse
// @Failure      401       {object} dto.ErrorResponse
// @Failure      404       {object} dto.ErrorResponse
// @Router       /bookings/{id} [get]
func (bc *BookingController) getBookingDetail(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    userRole, _ := ctx.Get("user_role")
    isAdmin := userRole == "admin"

    idParam := ctx.Param("id")
    bookingID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid booking ID"))
        return
    }

    response, err := bc.service.GetBookingDetail(bookingID, userID.(int), isAdmin)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// CancelBooking godoc
// @Summary      Batalkan booking
// @Description  Customer membatalkan booking yang masih aktif
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path     int                     true  "ID Booking"
// @Param        payload  body     dto.CancelBookingRequest true  "Alasan pembatalan"
// @Success      200      {object} dto.CancelBookingResponse
// @Failure      400      {object} dto.ErrorResponse
// @Failure      401      {object} dto.ErrorResponse
// @Router       /bookings/{id}/cancel [post]
func (bc *BookingController) cancelBooking(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    idParam := ctx.Param("id")
    bookingID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid booking ID"))
        return
    }

    var payload dto.CancelBookingRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := bc.service.CancelBooking(bookingID, userID.(int), payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// GetAllBookings godoc
// @Summary      Ambil semua booking (Admin Only)
// @Description  Mengambil semua booking dengan filter
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status   query string false "Filter status"
// @Param        page     query int    false "Halaman"
// @Param        limit    query int    false "Limit"
// @Success      200      {object} dto.BookingListResponse
// @Failure      400      {object} dto.ErrorResponse
// @Failure      401      {object} dto.ErrorResponse
// @Router       /bookings/admin [get]
func (bc *BookingController) getAllBookings(ctx *gin.Context) {
    var filter dto.BookingFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    response, err := bc.service.GetAllBookings(filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

// UpdateBookingStatus godoc
// @Summary      Update status booking (Admin Only)
// @Description  Admin mengubah status booking
// @Tags         Bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path     int                           true "ID Booking"
// @Param        payload  body     dto.UpdateBookingStatusRequest true "Status baru"
// @Success      200      {object} dto.UpdateBookingStatusResponse
// @Failure      400      {object} dto.ErrorResponse
// @Failure      401      {object} dto.ErrorResponse
// @Router       /bookings/admin/{id}/status [put]
func (bc *BookingController) updateBookingStatus(ctx *gin.Context) {
    idParam := ctx.Param("id")
    bookingID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid booking ID"))
        return
    }

    var payload dto.UpdateBookingStatusRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    response, err := bc.service.UpdateBookingStatus(bookingID, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, response)
}