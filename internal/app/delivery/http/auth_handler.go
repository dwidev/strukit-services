package http

import (
	"net/http"
	"strukit-services/internal/app/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	BaseController
}

func (a AuthHandler) LoginWithEmail(c *gin.Context) {
	var body dto.LoginWithEmailRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		logrus.Printf("login with email error binding request %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hit the services login with email
	c.JSON(200, gin.H{
		"message": "Login",
	})
}
