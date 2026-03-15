package authMiddleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zenrot/CryptoService/internal/auth"
)

func AuthMiddleware(auth auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		prefix := "Bearer "
		authToken := strings.TrimPrefix(authHeader, prefix)
		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Authorization header is empty")})
			c.Abort()
			return
		}
		err := auth.AuthorizeUser(authToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
