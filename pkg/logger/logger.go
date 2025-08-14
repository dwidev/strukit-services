package logger

import (
	"context"
	"os"
	"path/filepath"
	"strukit-services/pkg/constant"
	appContext "strukit-services/pkg/context"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Log *Logger

func New(cfg *Config) *Logger {

	logger := logrus.New()

	if cfg.Env == constant.Prod {
		logDir := "./.logs"
		os.MkdirAll(logDir, 0755)

		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetOutput(&lumberjack.Logger{
			Filename: filepath.Join(logDir, "strukit-prod.log"),
			MaxSize:  10,
			MaxAge:   30,
			Compress: true,
		})
	} else {
		// dev log
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logrus.SetReportCaller(true)

	Log = &Logger{
		Logger: logger,
	}
	return Log
}

type Config struct {
	Env constant.Environment
}

type Logger struct {
	*logrus.Logger
}

func (l *Logger) Request(ctx context.Context, body any) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)
	method := ctx.Value(appContext.MethodKey)
	path := ctx.Value(appContext.PathKey)
	ip := ctx.Value(appContext.IPAddressKey)

	return logrus.WithFields(logrus.Fields{
		"requestId": requestId,
		"method":    method,
		"path":      path,
		"user_ip":   ip,
		"body":      body,
	})
}

func (l *Logger) Handler(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return logrus.WithFields(logrus.Fields{
		"module":    "handler",
		"requestId": requestId,
	})
}

func (l *Logger) Service(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return logrus.WithFields(logrus.Fields{
		"module":    "services",
		"requestId": requestId,
	})
}

func (l *Logger) DB(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return logrus.WithFields(logrus.Fields{
		"module":    "DB",
		"requestId": requestId,
	})
}

func (l *Logger) LLM(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return logrus.WithFields(logrus.Fields{
		"module":    "llm-service",
		"requestId": requestId,
	})
}
