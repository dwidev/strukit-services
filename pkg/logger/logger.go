package logger

import (
	"context"
	appCtx "strukit-services/pkg/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (l *Logger) LogData(c *gin.Context, requestId string) logrus.Fields {

	return logrus.Fields{
		"requestId": requestId,
		"method":    c.Request.Method,
		"path":      c.Request.URL.Path,
	}
}

func (l *Logger) HttpRequestMiddlerware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Header("X-Request-ID", requestId)

		ctx := context.WithValue(c.Request.Context(), appCtx.RequestIDKey, requestId)
		c.Request = c.Request.WithContext(ctx)

		data := l.LogData(c, requestId)
		l.WithFields(data).Info("HTTP REQUEST LOG")
		c.Next()
	}
}
