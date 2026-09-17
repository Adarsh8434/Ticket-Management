package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/handler"
	"ticket-system/internal/middleware"
	"ticket-system/internal/repository"
	"ticket-system/internal/service"
)

func main() {

	// -------------------------
	// Database
	// -------------------------

	db, err := repository.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// -------------------------
	// User dependencies
	// -------------------------

	userRepository := repository.NewUserRepository(db)

	authService := service.NewAuthService(
		userRepository,
		"my-secret-key",
	)

	authHandler := handler.NewAuthHandler(authService)

	// -------------------------
	// Ticket dependencies
	// -------------------------

	ticketRepository := repository.NewTicketRepository(db)

	ticketService := service.NewTicketService(
		ticketRepository,
	)

	ticketHandler := handler.NewTicketHandler(
		ticketService,
	)

	// -------------------------
	// Router
	// -------------------------

	router := gin.Default()

	// Public endpoint
	//for checking health of the server
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Authentication
	//by router
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	// -------------------------
	// Protected ticket routes
	// -------------------------

	protected := router.Group("/tickets")

	protected.Use(
		middleware.AuthMiddleware("my-secret-key"),
	)

	protected.POST("", ticketHandler.Create)
	protected.GET("", ticketHandler.GetAll)
	protected.GET("/:id", ticketHandler.GetByID)
	protected.PATCH("/:id/status", ticketHandler.UpdateStatus)
	// -------------------------
	// Start server
	// -------------------------

	log.Println("Server starting on port 8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
