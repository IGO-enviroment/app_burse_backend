package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	zapConfig zap.Config
}

func NewLogger(options ...Option) *Logger {
	p := &Logger{
		zapConfig: zap.NewProductionConfig(),
	}

	p.zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	p.zapConfig.Development = false
	p.zapConfig.OutputPaths = []string{p.pwd() + "./logs/app.log"}

	for _, option := range options {
		option(p)
	}

	if p.zapConfig.Development {
		p.zapConfig.OutputPaths = append(p.zapConfig.OutputPaths, "stderr")
	}

	return p
}

func (l *Logger) Build() *zap.Logger {
	log, err := l.zapConfig.Build()
	if err != nil {
		log.Fatal("Failed to initialize logger")
	}

	return log
}

func (l *Logger) pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory")
	}

	return pwd
}
