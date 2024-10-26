package producer

import "go.uber.org/zap"

type Option func(*Producer)

func WithLogger(log *zap.Logger) Option {
	return func(p *Producer) {
		p.log = log
	}
}

func WithQueues(queues ...string) Option {
	return func(p *Producer) {
		p.queues = queues
	}
}

func WithTableName(tableName string) Option {
	return func(p *Producer) {
		p.tableName = tableName
	}
}
