package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

type Loggable struct {
	LogConfig zap.Config
	Log       *zap.Logger
}

func (log *Loggable) SetLogLevel(lvl zapcore.Level) {
	log.LogConfig.Level.SetLevel(lvl)
}

func New(name string) *Loggable {
	logConfig := zap.NewProductionConfig()
	logger, err := logConfig.Build()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &Loggable{
		LogConfig: logConfig,
		Log:       logger.Named(name),
	}
}
