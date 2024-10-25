package openai

import (
	"go.uber.org/zap"
)

type LeveledZapLogger struct {
	*zap.Logger
}

func (l *LeveledZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.Logger.Error(msg, fields(keysAndValues)...)
}

func (l *LeveledZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Info(msg, fields(keysAndValues)...)
}

func (l *LeveledZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.Logger.Debug(msg, fields(keysAndValues)...)
}

func (l *LeveledZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.Logger.Warn(msg, fields(keysAndValues)...)
}

func fields(keysAndValues []interface{}) []zap.Field {
	fields := []zap.Field{}
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		if v, ok := keysAndValues[i].(string); ok {
			fields = append(fields, zap.Any(v, keysAndValues[i+1]))
		}
	}

	return fields
}
