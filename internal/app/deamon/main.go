package deamon

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/jobs"
	"app_burse_backend/pkg/logger"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/consumer"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Deamon struct {
	config   *configs.Config
	log      *zap.Logger
	db       *sqlx.DB
	consumer *consumer.Consumer
}

func NewInstance(config *configs.Config) *Deamon {
	return &Deamon{
		config: config,
	}
}

func (d *Deamon) Setup() error {
	// Подключение к базе данных
	d.db = postgres.NewPostgres().Connect(d.config.DBHost, d.config.DBPort, "user", "pass", "postgres")

	// Загрузка конфигурации логгера
	d.log = logger.NewLogger(logger.WithDevelopment(true), logger.WithLevel(zap.DebugLevel)).Build()

	// Загрузка конфигурации очереди
	d.consumer = consumer.NewConsumer(
		consumer.WithDB(d.db),
		consumer.WithLogger(d.log),
	)
	jobs.RegisterJobs(d.consumer)

	return nil
}

func (d *Deamon) Run() error {
	ctx := context.Background()

	d.log.Info("Daemon started at", zap.Time("startup_time", time.Now()))

	go d.RunBackgroundTasks(ctx)
	go d.RunCronTasks(ctx)

	<-ctx.Done()

	d.log.Info("Daemon stopped at", zap.Time("shutdown_time", time.Now()))

	return nil
}

func (d *Deamon) RunCronTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Запуск крон задания
		}
	}
}

func (d *Deamon) RunBackgroundTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Запуск процесса выполнения фоновых заданий
			d.consumer.Run(ctx)
		}
	}
}
