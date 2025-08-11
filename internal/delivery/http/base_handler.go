package http

import (
	"strukit-services/pkg/logger"

	"gorm.io/gorm"
)

type BaseHandler struct {
	*gorm.DB
	*logger.Logger
}
