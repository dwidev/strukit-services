package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log *Logger

func New() *Logger {
	lg := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	Log = &Logger{
		Logger: lg,
	}
	return Log
}

type Logger struct {
	*logrus.Logger
}

func (l *Logger) LogRequest(c *gin.Context, requestId string) logrus.Fields {
	return logrus.Fields{
		"requestId": requestId,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
	}
}
