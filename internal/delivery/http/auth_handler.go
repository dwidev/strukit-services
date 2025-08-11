package http

import (
	"net/http"
	"strukit-services/internal/dto"
	"strukit-services/internal/services"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewAuth(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

type AuthHandler struct {
	BaseHandler
	*services.AuthService
}

func (a AuthHandler) LoginWithEmail(c *gin.Context) {
	var body dto.LoginWithEmailRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Log.Errorf("error binding request %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := a.AuthService.LoginWithEmail(body.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}
