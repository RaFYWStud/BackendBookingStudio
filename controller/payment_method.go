package controller

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/RaFYWStud/BackendBookingStudio/config/middleware"
    "github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
    "github.com/RaFYWStud/BackendBookingStudio/contract"
    "github.com/RaFYWStud/BackendBookingStudio/dto"
)

type PaymentMethodController struct {
    service contract.PaymentMethodService
}

func (pmc *PaymentMethodController) GetPrefix() string {
    return "/payment-methods"
}

func (pmc *PaymentMethodController) InitService(service *contract.Service) {
    pmc.service = service.PaymentMethod
}

func (pmc *PaymentMethodController) InitRoute(app *gin.RouterGroup) {
    // Public routes - Customer view active payment methods
    app.GET("", pmc.getActivePaymentMethods)

    // Admin routes - Manage payment methods
    admin := app.Group("/admin")
    admin.Use(middleware.Auth(), middleware.AdminOnly())
    {
        admin.GET("", pmc.getAllPaymentMethods)
        admin.GET("/:id", pmc.getPaymentMethodByID)
        admin.POST("", pmc.createPaymentMethod)
        admin.PUT("/:id", pmc.updatePaymentMethod)
        admin.DELETE("/:id", pmc.deletePaymentMethod)
    }
}

// GetActivePaymentMethods - GET /payment-methods (Public)
// Get all active payment methods that customers can use
func (pmc *PaymentMethodController) getActivePaymentMethods(ctx *gin.Context) {
    result, err := pmc.service.GetActivePaymentMethods()
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// GetAllPaymentMethods - GET /payment-methods/admin (Admin)
// Get all payment methods (active & inactive)
func (pmc *PaymentMethodController) getAllPaymentMethods(ctx *gin.Context) {
    result, err := pmc.service.GetAllPaymentMethods()
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// GetPaymentMethodByID - GET /payment-methods/admin/:id (Admin)
// Get payment method detail by ID
func (pmc *PaymentMethodController) getPaymentMethodByID(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid payment method ID"))
        return
    }

    result, err := pmc.service.GetPaymentMethodByID(id)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// CreatePaymentMethod - POST /payment-methods/admin (Admin)
// Create new payment method
func (pmc *PaymentMethodController) createPaymentMethod(ctx *gin.Context) {
    var payload dto.CreatePaymentMethodRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    result, err := pmc.service.CreatePaymentMethod(payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusCreated, result)
}

// UpdatePaymentMethod - PUT /payment-methods/admin/:id (Admin)
// Update existing payment method
func (pmc *PaymentMethodController) updatePaymentMethod(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid payment method ID"))
        return
    }

    var payload dto.UpdatePaymentMethodRequest
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        HandlerError(ctx, errs.BadRequest("invalid request payload"))
        return
    }

    result, err := pmc.service.UpdatePaymentMethod(id, payload)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// DeletePaymentMethod - DELETE /payment-methods/admin/:id (Admin)
// Delete payment method
func (pmc *PaymentMethodController) deletePaymentMethod(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        HandlerError(ctx, errs.BadRequest("invalid payment method ID"))
        return
    }

    result, err := pmc.service.DeletePaymentMethod(id)
    if err != nil {
        HandlerError(ctx, err)
        return
    }

    ctx.JSON(http.StatusOK, result)
}