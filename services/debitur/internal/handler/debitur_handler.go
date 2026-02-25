package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-debitur/internal/service"
)

type DebiturHandler struct {
	Service *service.DebiturService
}

func (h *DebiturHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "debitur service is running"})
}
