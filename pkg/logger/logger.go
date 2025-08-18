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
		// prod log
		logDir := "./.logs"
		os.MkdirAll(logDir, 0755)

		logger.SetLevel(logrus.InfoLevel)
		logger.SetOutput(&lumberjack.Logger{
			Filename: filepath.Join(logDir, "strukit-prod.log"),
			MaxSize:  10,
			MaxAge:   30,
			Compress: true,
		})
	} else {
		// dev log
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(os.Stdout)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logger.SetReportCaller(true)

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

	return l.WithFields(logrus.Fields{
		"requestId": requestId,
		"method":    method,
		"path":      path,
		"user_ip":   ip,
		"body":      body,
	})
}

func (l *Logger) Handler(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return l.WithFields(logrus.Fields{
		"module":    "handler",
		"requestId": requestId,
	})
}

func (l *Logger) Service(ctx context.Context, data ...any) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)
	f := logrus.Fields{
		"module":    "services",
		"requestId": requestId,
	}
	if len(data) > 0 {
		f["req_data"] = data
	}
	return l.WithFields(f)
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
	userId := ctx.Value(appContext.UserIDKey)

	return l.WithFields(logrus.Fields{
		"module":    "llm-service",
		"userId":    userId,
		"requestId": requestId,
	})
}

func (l *Logger) App(ctx context.Context) *logrus.Entry {
	requestId := ctx.Value(appContext.RequestIDKey)

	return l.WithFields(logrus.Fields{
		"module":    "app-error",
		"requestId": requestId,
	})
}
