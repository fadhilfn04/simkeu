package handler

import (
	"github.com/gin-gonic/gin"
	"simkeu/service-debitur/internal/service"
)

type DebiturHandler struct {
	Service *service.DebiturService
}

func (h *DebiturHandler) GetStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "debitur service is running",
	})
}

func (h *DebiturHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	data, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "debitur not found"})
		return
	}

	c.JSON(200, data)
}

func (h *DebiturHandler) Create(c *gin.Context) {
	var input struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	err := h.Service.Create(input.ID, input.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create debitur"})
		return
	}

	c.JSON(201, gin.H{"message": "debitur created"})
}
