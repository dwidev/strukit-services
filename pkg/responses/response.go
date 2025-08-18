package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

func New(StatusCode int, Message any) *MessageResponse {
	return &MessageResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
}

type MessageResponse struct {
	StatusCode int `json:"statusCode"`
	Message    any `json:"message"`
}

func ServerError(ctx *gin.Context) {
	code := http.StatusInternalServerError
	ctx.JSON(code, MessageResponse{StatusCode: code, Message: "There was an error on the server, please try again later."})
}
