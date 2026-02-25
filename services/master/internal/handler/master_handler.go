package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-master/internal/service"
)

type MasterHandler struct {
	Service *service.MasterService
}

func (h *MasterHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "master service is running"})
}
