package http

import (
	"strukit-services/pkg/logger"

	"gorm.io/gorm"
)

type BaseController struct {
	*gorm.DB
	*logger.Logger
}
