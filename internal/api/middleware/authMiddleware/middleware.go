package authMiddleware

import (
	"CryptoService/internal/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(auth auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		prefix := "Bearer "
		authToken := strings.TrimPrefix(authHeader, prefix)
		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Authorization header is empty.")})
			c.Abort()
			return
		}
		_, err := auth.AuthorizeUser(authToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
