package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(minRank int) gin.HandlerFunc {
	return func(c *gin.Context) {

		rankValue, exists := c.Get("rank")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Authorization data missing."})
			c.Abort()
			return 
		}

		userRank, ok := rankValue.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization format."})
			c.Abort()
			return 
		}

		if userRank < minRank {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions.",
			})
			c.Abort()
			return 
		}

		c.Next()
	}
}