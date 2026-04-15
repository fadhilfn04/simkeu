package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"simkeu/service-debitur/internal/database"
	"simkeu/service-debitur/internal/handler"
	"simkeu/service-debitur/internal/repository"
	"simkeu/service-debitur/internal/service"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	authMiddleware := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])

		c.Next()
	}

	// =====================
	// Database Connection
	// =====================
	db := database.Connect()
	defer db.Close()

	// Create tables if they don't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS debitur (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create debitur table:", err)
	}

	log.Println("Debitur table ready.")

	// =====================
	// Initialize Components
	// =====================
	debiturRepo := &repository.DebiturRepository{DB: db}
	debiturService := &service.DebiturService{Repo: debiturRepo}
	debiturHandler := &handler.DebiturHandler{Service: debiturService}

	// =====================
	// HTTP Server
	// =====================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Debitur service is healthy"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/status", debiturHandler.GetStatus)
		protected.GET("/debitur/:id", debiturHandler.GetByID)
		protected.POST("/debitur", debiturHandler.Create)
	}

	log.Printf("Debitur service running on port %s\n", port)
	router.Run(":" + port)
}
