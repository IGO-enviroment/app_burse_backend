package consumer

import (
	"app_burse_backend/pkg/queue/job"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Consumer struct {
	queues    []string
	tableName string

	handlers map[string]func(*job.ProcessJob) error

	db  *sqlx.DB
	log *zap.Logger
}

func NewConsumer(options ...Options) *Consumer {
	c := &Consumer{
		queues:    []string{"default"},
		tableName: "queue_jobs",
		handlers:  make(map[string]func(*job.ProcessJob) error),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Consumer) RegisterHandler(handlerName string, handler func(*job.ProcessJob) error) {
	// Регистрация обработчика
	c.handlers[handlerName] = handler
}

func (c *Consumer) Run(ctx context.Context) {
	for _, queueName := range c.queues {
		go c.workersQueue(ctx, queueName)
	}
	<-ctx.Done()
}

func (c *Consumer) workersQueue(ctx context.Context, queueName string) {
	workers := 10
	wg := &sync.WaitGroup{}

	defer wg.Wait()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					item, err := c.getItem(ctx, queueName)
					if item == nil {
						continue
					}

					if err != nil {
						c.log.Log(zap.ErrorLevel, "Ошибка получения из очереди", zap.Error(err))
						continue
					}

					c.log.Info("Обработка из очереди", zap.String("queue_name", queueName), zap.Int("job_id", item.ID))
					job := job.NewProcessJob(item, c.db, c.log)
					jobErr := job.Start(&c.handlers)

					if jobErr != nil {
						c.log.Log(zap.ErrorLevel, "Ошибка обработки из очереди", zap.Error(err))
						continue
					}

					c.ProcessedItem(item.ID)
				}
			}
		}(wg)
	}
}

func (c *Consumer) getItem(ctx context.Context, queueName string) (*job.RawJob, error) {
	var item job.RawJob

	// , processed = true
	err := c.db.QueryRowx(
		fmt.Sprintf(
			`UPDATE %s
			SET reserv_at = $1
			WHERE id IN (
			  SELECT id FROM %s
				WHERE 
					processed = false AND
					queue_name = $2 AND
					((run_at < $3 OR run_at IS NULL) AND (reserv_at IS NULL OR reserv_at < $4))
				ORDER BY id ASC LIMIT 1
			)
			RETURNING id, method, item`,
			c.tableName, c.tableName,
		),
		time.Now().UTC().Add(time.Second*(5*60)),
		queueName,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&item.ID, &item.Method, &item.Payload)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		c.log.Log(zap.ErrorLevel, "Ошибка получения из очереди", zap.Error(err))
		return nil, err
	}

	return &item, nil
}

func (c *Consumer) ProcessedItem(itemID int) {
	_, err := c.db.Exec(fmt.Sprintf(`UPDATE %s SET processed = true WHERE id = $1`, c.tableName), itemID)
	if err != nil {
		c.log.Log(zap.ErrorLevel, "Ошибка обработки из очереди", zap.Error(err))
	}
}
