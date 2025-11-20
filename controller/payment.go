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

type PaymentController struct {
    service contract.PaymentService
}

func (pc *PaymentController) GetPrefix() string {
    return "/payments"
}

func (pc *PaymentController) InitService(service *contract.Service) {
    pc.service = service.Payment
}

func (pc *PaymentController) InitRoute(app *gin.RouterGroup) {
    // Customer routes (require authentication)
    customer := app.Group("")
    customer.Use(middleware.Auth())
    {
        customer.POST("/bookings/:booking_id/upload", pc.uploadPaymentProof)
        customer.GET("", pc.getMyPayments)
        customer.GET("/:id", pc.getPaymentDetail)
    }

    // Admin routes
    admin := app.Group("/admin")
    admin.Use(middleware.Auth(), middleware.AdminOnly())
    {
        admin.GET("", pc.getAllPayments)
        admin.GET("/pending", pc.getPendingPayments)
        admin.POST("/:id/verify", pc.verifyPayment)
    }
}

// UploadPaymentProof - POST /payments/bookings/:booking_id/upload (Customer)
// Upload payment proof for a booking (DP or full payment)
func (pc *PaymentController) uploadPaymentProof(ctx *gin.Context) {
    // Get user ID from JWT middleware
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    // Get booking ID from URL parameter
    bookingIDParam := ctx.Param("booking_id")
    bookingID, err := strconv.Atoi(bookingIDParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid booking ID"))
        return
    }

    // Bind request payload
    var payload dto.UploadPaymentProofRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    // Call service
    result, err := pc.service.UploadPaymentProof(userID.(int), bookingID, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, result)
}

// GetMyPayments - GET /payments (Customer)
// Get payment history for logged-in customer
// Query params: status, type, booking_id, page, limit
func (pc *PaymentController) getMyPayments(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    // Parse query parameters
    var filter dto.PaymentFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    // Set defaults if not provided
    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    result, err := pc.service.GetMyPayments(userID.(int), filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// GetPaymentDetail - GET /payments/:id (Customer/Admin)
// Get payment detail by ID
func (pc *PaymentController) getPaymentDetail(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    idParam := ctx.Param("id")
    paymentID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid payment ID"))
        return
    }

    // Check if user is admin
    userRole, _ := ctx.Get("user_role")
    isAdmin := userRole == "admin"

    result, err := pc.service.GetPaymentDetail(paymentID, userID.(int), isAdmin)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// GetAllPayments - GET /payments/admin (Admin)
// Get all payments with filters
// Query params: booking_id, status, type, page, limit
func (pc *PaymentController) getAllPayments(ctx *gin.Context) {
    var filter dto.PaymentFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    result, err := pc.service.GetAllPayments(filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// GetPendingPayments - GET /payments/admin/pending (Admin)
// Get all pending payments waiting for verification
// Query params: page, limit
func (pc *PaymentController) getPendingPayments(ctx *gin.Context) {
    var filter dto.PaymentFilterRequest
    if err := ctx.ShouldBindQuery(&filter); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid query parameters"))
        return
    }

    if filter.Page <= 0 {
        filter.Page = 1
    }
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    result, err := pc.service.GetPendingPayments(filter)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// VerifyPayment - POST /payments/admin/:id/verify (Admin)
// Verify or reject a payment
func (pc *PaymentController) verifyPayment(ctx *gin.Context) {
    adminID, exists := ctx.Get("user_id")
    if !exists {
        HandlerError(ctx, errs.Unauthorized("user not authenticated"))
        return
    }

    idParam := ctx.Param("id")
    paymentID, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid payment ID"))
        return
    }

    var payload dto.VerifyPaymentRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    result, err := pc.service.VerifyPayment(paymentID, adminID.(int), payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}