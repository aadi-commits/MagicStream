package middlewares

import (
	"net/http"
	"strings"

	"github.com/aadi-commits/MagicStream/Server/MagicStreamServer/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing."})
			c.Abort()
			return 
		}

		//Excep: Bearer (token)
		tokenString :=  strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format."})
			c.Abort()
			return 
		}

		claims, err := utils.ValidateAccessToken(tokenString) 
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token."})
			c.Abort()
			return 
		}

		rank, exists := utils.RoleRank[claims.Role]
		if !exists{
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role."})
			c.Abort()
			return 
		}

		// Attach user information
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("rank", rank)

		c.Next()

	}
}