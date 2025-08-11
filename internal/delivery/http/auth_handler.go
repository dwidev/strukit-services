package http

import (
	"net/http"
	"strukit-services/internal/dto"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	BaseHandler
}

func (a AuthHandler) LoginWithEmail(c *gin.Context) {
	var body dto.LoginWithEmailRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Log.Errorf("error binding request %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hit the services login with email
	c.JSON(200, gin.H{
		"message": "Login",
	})
}
