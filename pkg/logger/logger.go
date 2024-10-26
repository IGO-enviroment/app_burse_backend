package logger

type Logger struct {
}

func NewLogger(options ...Option) *Logger {
	p := &Logger{}

	for _, option := range options {
		option(p)
	}

	return p
}

func Setup() error {
	// Setup logger
	return nil
}
