package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-tagihan/internal/service"
)

type TagihanHandler struct {
	Service *service.TagihanService
}

func (h *TagihanHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "tagihan service is running"})
}
