package responses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	StatusCode int `json:"statusCode"`
	Data       any `json:"data"`
}

func New(StatusCode int, Message any) MessageResponse {
	return MessageResponse{
		StatusCode: StatusCode,
		Message:    Message,
	}
}

func Created(resource ...string) MessageResponse {
	msg := "Resource has been successfully created"
	if len(resource) > 0 {
		msg = fmt.Sprintf("%s has been successfully created", resource[0])
	}

	return MessageResponse{
		StatusCode: http.StatusCreated,
		Message:    msg,
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
