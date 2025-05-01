package logger

import (
	"rapid-bridge/domain/port"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

// CustomTimeEncoder formats time in a human-readable format.
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339)) // You can change the format if needed
}

func NewZapLogger() (port.Logger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = CustomTimeEncoder
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return &zapLogger{
		logger: sugar,
	}, nil
}

func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debugw(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.logger.Infow(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warnw(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.logger.Errorw(msg, fields...)
}

func (l *zapLogger) Panic(msg string, fields ...interface{}) {
	l.logger.Panicw(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatalw(msg, fields...)
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
