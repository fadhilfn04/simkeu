package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"simkeu/service-auth/internal/database"
	"simkeu/service-auth/internal/handler"
	"simkeu/service-auth/internal/repository"
	"simkeu/service-auth/internal/service"
)

func main() {

	db := database.Connect()

	userRepo := &repository.UserRepository{DB: db}
	authService := &service.AuthService{
		Repo:       userRepo,
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DebiturURL: os.Getenv("DEBITUR_SERVICE_URL"),
	}
	authHandler := &handler.AuthHandler{
		Service: authService,
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/validate", authHandler.Validate)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}