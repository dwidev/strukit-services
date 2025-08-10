package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log *Logger

func New() *Logger {
	logrus := logrus.New()
	Instance := &Logger{
		Logger: logrus,
	}
	return Instance
}

type Logger struct {
	*logrus.Logger
}

func (l *Logger) HttpRequestMiddlerware() gin.HandlerFunc {
	return func(c *gin.Context) {
		l.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Http request")

		c.Next()
	}
}
