package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/RaFYWStud/BackendBookingStudio/config/middleware"
	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	// GetPrefix returns the route prefix that the controller will use.
	GetPrefix() string

	// InitService initializes the necessary services for the controller.
	// This service typically contains the business logic required by the controller.
	InitService(service *contract.Service)

	// InitRoute sets up the routes for the controller within the given router group.
	InitRoute(app *gin.RouterGroup)
}

func New(app *gin.Engine, service *contract.Service) {
	allController := []Controller{
		&AuthController{},
		&StudioController{},
		&BookingController{},
		&PaymentMethodController{},
		&PaymentController{},
		// Add your controller here
	}

	// do not modify the code below there
	for _, c := range allController {
		c.InitService(service)
		group := app.Group(c.GetPrefix())
		group.Use(middleware.CORSMiddleware())
		c.InitRoute(group)
		log.Printf("initiate route %s\n", c.GetPrefix())
	}
}

// handlerError is a helper function to handle errors in the controller.
// It checks if the error is of type MessageError and responds with the appropriate status code and message.
func HandlerError(ctx *gin.Context, err error) {
	var messageErr errs.MessageError
	if errors.As(err, &messageErr) {
		ctx.JSON(messageErr.Status(), messageErr)
		return
	}
	_ = ctx.Error(err).SetType(gin.ErrorTypePrivate) // record internal error
	ctx.JSON(http.StatusInternalServerError, errs.InternalServerError("Internal Server Error"))
}
