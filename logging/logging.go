package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hemagome/hefesto/config"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("configuration not loaded")
	}

	// Create logs directory if it doesn't exist
	logsDir := cfg.Logging.Path
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}

	// Configure the log file path with current month
	currentTime := time.Now()
	logFile := filepath.Join(logsDir, fmt.Sprintf("hefesto-%s.json", currentTime.Format("2006-01")))

	// Create or open the log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	// Create encoder configuration for JSON format
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core with JSON encoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	// Create logger
	Logger = zap.New(core, zap.AddCaller())
}

func SetLogLevel(level string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// Create a new atomic level
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	// Update the logger's level
	Logger = Logger.WithOptions(zap.IncreaseLevel(atomicLevel))
}
