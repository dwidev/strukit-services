package http

import (
	"net/http"
	"strukit-services/internal/dto"
	"strukit-services/internal/services"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
)

func NewAuth(base *BaseHandler, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: base,
		AuthService: authService,
	}
}

type AuthHandler struct {
	*BaseHandler
	*services.AuthService
}

func (a AuthHandler) CreatePassword(c *gin.Context) {
	var body dto.CreatePasswordDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Log.Errorf("error binding request %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if msg := a.AppValidator.Valid(&body); len(msg) > 0 {
		err := responses.BodyErr(msg)
		c.Error(err)
		return
	}

	err := a.AuthService.CreatePassword(c.Request.Context(), body.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, responses.Created("Password"))
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

	c.JSON(http.StatusOK, response)
}
