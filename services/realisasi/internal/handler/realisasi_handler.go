package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-realisasi/internal/service"
)

type RealisasiHandler struct {
	Service *service.RealisasiService
}

func (h *RealisasiHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "realisasi service is running"})
}
