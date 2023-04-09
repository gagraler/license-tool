package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

// LogLevel defines the severity of a log message.
type LogLevel uint8

const (
	// DebugLevel logs debug messages.
	DebugLevel LogLevel = iota
	// InfoLevel logs informational messages.
	InfoLevel
	// WarnLevel logs warning messages.
	WarnLevel
	// ErrorLevel logs error messages.
	ErrorLevel
	// FatalLevel logs critical messages and then calls os.Exit(1).
	FatalLevel
)

// Config holds logger configuration options.
type Config struct {
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LogLevel   LogLevel
}

// EncoderConfig holds encoder configuration options.
type EncoderConfig struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string
	LineEnding    string
	EncodeLevel   zapcore.LevelEncoder
	EncodeTime    zapcore.TimeEncoder
	EncodeCaller  zapcore.CallerEncoder
}

// InitLogger initializes the default logger with the given configuration.
func InitLogger(config Config) error {
	hook := lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
	write := zapcore.AddSync(&hook)
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(config.LogLevel))
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, write, atomicLevel)
	caller := zap.AddCaller()
	development := zap.Development()
	defaultLogger = zap.New(core, caller, development)
	return nil
}

// Debug logs a debug message using the default logger.
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

// Info logs an informational message using the default logger.
func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error logs an error message using the default logger.
func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

// Fatal logs a critical message and then calls os.Exit(1) using the default logger.
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

// WithFields returns a slice of fields given a list of key-value pairs.
func WithFields(data map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(config Config) (*zap.Logger, error) {
	hook := lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
	writeSyncer := zapcore.AddSync(&hook)
	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(config.LogLevel))
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)

	logger := zap.New(core, zap.AddCaller(), zap.Development())
	return logger, nil
}
