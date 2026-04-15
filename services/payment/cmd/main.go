package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"simkeu/service-payment/internal/database"
	"simkeu/service-payment/internal/handler"
	"simkeu/service-payment/internal/repository"
	"simkeu/service-payment/internal/service"
	"simkeu/service-payment/internal/middleware"
)

func main() {

	// =====================
	// Database Connection
	// =====================
	db := database.Connect()
	defer db.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		amount DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create payments table:", err)
	}

	log.Println("Payments table ready.")

	// =====================
	// Initialize Components
	// =====================
	paymentRepo := &repository.PaymentRepository{DB: db}
	paymentService := &service.PaymentService{Repo: paymentRepo}
	paymentHandler := &handler.PaymentHandler{Service: paymentService}

	// =====================
	// HTTP Server
	// =====================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Payment service is healthy"})
	})

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.GET("/status", paymentHandler.GetStatus)
	}

	log.Printf("Payment service running on port %s\n", port)
	router.Run(":" + port)
}