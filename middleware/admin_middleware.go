package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.MustGet("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		role := claims["role"]

		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin privileges required",
			})
			return
		}

		c.Next()
	}
}
