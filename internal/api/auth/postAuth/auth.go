package postAuth

import (
	"CryptoService/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type requestPostAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(auth auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requestPostAuth
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokenStr, err := auth.AuthenticateUser(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	}
}

func RegisterHandler(auth auth.Authorizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requestPostAuth
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokenStr, err := auth.RegisterUser(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	}
}
