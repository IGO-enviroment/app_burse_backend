package postgres

import (
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
func (p *Postgres) Connect(host string, port int, user, password, dbname string) *sqlx.DB {
	var db *sqlx.DB
	var err error

	for i := 0; i < p.connAttempts; i++ {
		db, err = sqlx.Connect("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
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
