// file: logger.go
package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Logger() (*zap.Logger, error) {

	// 创建 lumberjack 日志切割器
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    logMaxSize, // MB
		MaxBackups: logMaxBackups,
		MaxAge:     logMaxAge, // Days
		Compress:   true,
	}

	// 设置日志输出格式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// 选择输出 txt 格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 设置堆栈跟踪
	stacktraceLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	stacktrace := zap.AddCaller(encoder)

	// 配置核心日志写入
	core := zapcore.NewTee(
		// 控制台输出
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zap.InfoLevel),

		// 文件输出
		zapcore.NewCore(encoder, zapcore.AddSync(lumberjackLogger), zap.InfoLevel),

		// 堆栈跟踪
		zapcore.NewCore(stacktrace, zapcore.AddSync(lumberjackLogger), stacktraceLevel),
	)

	// 构建日志
	logger := zap.New(core)

	// 保存日志实例，以便在其他地方调用
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)

	// 添加 http 请求和响应日志
	//httpLogger := logger.With(
	//	zap.Namespace("http"),
	//	zap.String("protocol", "http"),
	//)
	return logger, nil
}
