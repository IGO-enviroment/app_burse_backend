package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*Logger)

func WithLevel(level zapcore.Level) Option {
	return func(l *Logger) {
		l.zapConfig.Level = zap.NewAtomicLevelAt(level)
	}
}

func WithDevelopment(development bool) Option {
	return func(l *Logger) {
		l.zapConfig.Development = development
	}
}

func WithOutputPaths(outputPaths ...string) Option {
	return func(l *Logger) {
		l.zapConfig.OutputPaths = outputPaths
	}
}
