package responses

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AppForbidden      = &AppError{http.StatusForbidden, "Access denied: You don't own this project"}
	DateLessThanToday = &AppError{http.StatusBadRequest, "date/time is less than today"}
)

func Err(Code int, Message any) *AppError {
	return &AppError{Code, Message}
}

func BodyErr(Message any) *AppError {
	return &AppError{http.StatusBadRequest, Message}
}

type AppError struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

func (e AppError) Error() string {
	if message, ok := e.Message.(string); ok {
		return message
	}

	if slice, ok := e.Message.([]string); ok {
		message := strings.Join(slice, ", ")
		return message
	}

	return fmt.Sprintf("%v", e.Message)
}

func (e *AppError) JSON(ctx *gin.Context) {
	ctx.JSON(e.Code, MessageResponse{StatusCode: e.Code, Message: e.Message})
}
