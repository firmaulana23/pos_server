package app

import (
	"fmt"
	"log"
	"pos-system/internal/config"
	"pos-system/internal/database"
	"pos-system/internal/routes"
	"pos-system/pkg/auth"

	"github.com/gin-gonic/gin"
)

type App struct {
	config     *config.Config
	database   *database.Database
	jwtService *auth.JWTService
	router     *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWT.SecretKey, cfg.JWT.ExpiryHours)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db.DB, jwtService)

	return &App{
		config:     cfg,
		database:   db,
		jwtService: jwtService,
		router:     router,
	}
}

func (a *App) Run() error {
	address := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	log.Printf("Starting server on %s", address)
	return a.router.Run(address)
}
