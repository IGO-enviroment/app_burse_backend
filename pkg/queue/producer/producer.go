package producer

import (
	"app_burse_backend/pkg/queue/job"
	"encoding/json"
	"fmt"
	"time"

	"slices"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const defaultQueue = "default"
const defaultTableName = "queue_jobs"

type Producer struct {
	log    *zap.Logger
	queues []string

	tableName string
}

func NewProducerQueue(options ...Option) *Producer {
	queues := []string{defaultQueue}
	p := &Producer{
		queues:    queues,
		log:       zap.NewNop(),
		tableName: defaultTableName,
	}

	for _, opt := range options {
		opt(p)
	}

	return p
}

type AddOptions struct {
	DB *sqlx.DB

	QueueName string
	RunAt     time.Time

	Job job.RawJob
}

func (q *Producer) Add(options AddOptions) error {
	if !slices.Contains(q.queues, options.QueueName) {
		q.log.Log(zap.ErrorLevel, "Отсутствует такая очередь", zap.String("queue_name", options.QueueName))
		return fmt.Errorf("отсутствует такая очередь: %s", options.QueueName)
	}

	jsonBytes, err := json.MarshalIndent(options.Job.Params, "", "  ")
	if err != nil {
		q.log.Log(zap.ErrorLevel, "Ошибка сериализации данных", zap.Error(err))
		return err
	}

	// Add item to the queue
	_, err = options.DB.Exec(
		fmt.Sprintf("INSERT INTO %s (queue_name, method, item, run_at) VALUES ($1, $2, $3, $4)", q.tableName),
		options.QueueName,
		options.Job.Method,
		string(jsonBytes),
		options.RunAt.UTC(),
	)
	if err != nil {
		q.log.Log(zap.ErrorLevel, "Ошибка добавления в очередь", zap.Error(err))
		return err
	}

	q.log.Log(zap.InfoLevel, "Записано в очередь", zap.String("queue_name", options.QueueName))
	return nil
}
