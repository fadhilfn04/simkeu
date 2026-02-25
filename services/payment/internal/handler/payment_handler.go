package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-payment/internal/service"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func (h *PaymentHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "payment service is running"})
}
