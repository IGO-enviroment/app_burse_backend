package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

const (
	defaultMaxPoolSize  = 10
	defaultConnAttempts = 10
	defaultConnTimeout  = 5
)

// Интерфейс для работы с базой данных.
type Database interface {
	sqlx.Preparer
	sqlx.PreparerContext
	sqlx.Ext
	sqlx.ExtContext

	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Close() error
}

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

func NewPostgres(options ...Option) *Postgres {
	p := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout * time.Second,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// Инициализация подключения.
func (p *Postgres) Connect(host string, port int, user, password, dbname string) Database {
	var db *sqlx.DB
	var err error

	for i := 0; i < p.connAttempts; i++ {
		db, err = sqlx.Connect("postgres", fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, password, host, port, dbname,
		))

		if err == nil {
			break
		}

		time.Sleep(p.connTimeout)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
