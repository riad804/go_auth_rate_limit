package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/riad804/go_auth/internal/utils"
)

func AuthMiddleware(jwtUtil *utils.JWTWrapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		claims, err := jwtUtil.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("org_id", claims.OrgID)
		c.Next()
	}
}
