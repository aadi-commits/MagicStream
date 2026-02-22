package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/internal/ai"
	"github.com/gin-gonic/gin"
)

type AIController struct {
	service *ai.Service
}

func InitAIController(service *ai.Service) *AIController {
	return &AIController{service: service}
}

func (ctl *AIController) Infer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Prompt string `json:"prompt"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		prompt := strings.TrimSpace(input.Prompt)
		if prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt cannot be empty"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
		defer cancel()

		reply, err := ctl.service.GenerateReply(ctx, prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"reply": reply})
	}
}