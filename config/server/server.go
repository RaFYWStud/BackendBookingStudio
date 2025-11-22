package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/config"
	"github.com/RaFYWStud/BackendBookingStudio/config/database"
	"github.com/RaFYWStud/BackendBookingStudio/config/middleware"
	"github.com/RaFYWStud/BackendBookingStudio/controller"
	"github.com/gin-gonic/gin"

	"github.com/RaFYWStud/BackendBookingStudio/repository"
	"github.com/RaFYWStud/BackendBookingStudio/service"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/RaFYWStud/BackendBookingStudio/docs"
)

func Run() {
	log.Println("Starting application...")

	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
		return
	}

	db, _, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return
	}

	// Default: start the server
	startServer(cfg, db)
}

func startServer(cfg *config.AppConfigurationMap, db *gorm.DB) {
	// Initialize repositories and services
	repo := repository.New(db)
	serv := service.New(repo)

	// Set Gin mode
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.GlobalRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst, map[string]struct{}{}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Static("/static", "./static")

	// route Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes
	controller.New(r, serv)

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server is running on port %d", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
