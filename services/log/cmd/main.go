package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"simkeu/service-log/internal/database"
	"simkeu/service-log/internal/handler"
	"simkeu/service-log/internal/repository"
	"simkeu/service-log/internal/service"
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
	CREATE TABLE IF NOT EXISTS logs (
		id SERIAL PRIMARY KEY,
		action TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create logs table:", err)
	}

	log.Println("Logs table ready.")

	// =====================
	// Initialize Components
	// =====================
	logRepo := &repository.LogRepository{DB: db}
	logService := &service.LogService{Repo: logRepo}
	logHandler := &handler.LogHandler{Service: logService}

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
		c.JSON(200, gin.H{"message": "Log service is healthy"})
	})

	// Protected routes
	protected := router.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/status", logHandler.GetStatus)
	}

	log.Printf("Log service running on port %s\n", port)
	router.Run(":" + port)
}
