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

// CreateBooking - POST /bookings (Customer)
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

// GetMyBookings - GET /bookings (Customer)
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

// GetBookingDetail - GET /bookings/:id (Customer/Admin)
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

// CancelBooking - POST /bookings/:id/cancel (Customer)
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

// GetAllBookings - GET /bookings/admin (Admin)
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

// UpdateBookingStatus - PUT /bookings/admin/:id/status (Admin)
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