package consumer

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Options func(*Consumer)

func WithDB(db *sqlx.DB) Options {
	return func(c *Consumer) {
		c.db = db
	}
}

func WithLogger(log *zap.Logger) Options {
	return func(c *Consumer) {
		c.log = log
	}
}
