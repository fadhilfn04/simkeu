package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"simkeu/service-master/internal/database"
	"simkeu/service-master/internal/handler"
	"simkeu/service-master/internal/repository"
	"simkeu/service-master/internal/service"
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
	CREATE TABLE IF NOT EXISTS master (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create master table:", err)
	}

	log.Println("Master table ready.")

	// =====================
	// Initialize Components
	// =====================
	masterRepo := &repository.MasterRepository{DB: db}
	masterService := &service.MasterService{Repo: masterRepo}
	masterHandler := &handler.MasterHandler{Service: masterService}

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
		c.JSON(200, gin.H{"message": "Master service is healthy"})
	})

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/status", masterHandler.GetStatus)
	}

	log.Printf("Master service running on port %s\n", port)
	router.Run(":" + port)
}
