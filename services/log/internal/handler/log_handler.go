package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-log/internal/service"
)

type LogHandler struct {
	Service *service.LogService
}

func (h *LogHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "log service is running"})
}
