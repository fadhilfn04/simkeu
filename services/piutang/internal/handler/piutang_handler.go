package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-piutang/internal/service"
)

type PiutangHandler struct {
	Service *service.PiutangService
}

func (h *PiutangHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "piutang service is running"})
}
