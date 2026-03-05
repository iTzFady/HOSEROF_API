package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {

		value, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		claims := value.(*Claims)

		if claims.Role != "teacher" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "teacher privileges required",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		c.Next()
	}
}
