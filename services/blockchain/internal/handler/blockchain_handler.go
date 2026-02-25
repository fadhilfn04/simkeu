package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simkeu/service-blockchain/internal/service"
)

type BlockchainHandler struct {
	Service *service.BlockchainService
}

func (h *BlockchainHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "blockchain service is running"})
}
