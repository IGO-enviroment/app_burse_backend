package deamon

import (
	"app_burse_backend/configs"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/consumer"
	"context"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Deamon struct {
	config   *configs.Config
	log      zap.Logger
	db       *sqlx.DB
	consumer *consumer.Consumer
}

func NewInstance(config *configs.Config) *Deamon {
	return &Deamon{
		config: config,
	}
}

func (d *Deamon) Setup() error {
	d.db = postgres.NewPostgres().Connect(d.config.DBHost, d.config.DBPort, "user", "pass", "postgres")

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory")
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConfig.Development = true
	logConfig.OutputPaths = []string{"stdout", pwd + "/logs/app.log"}
	logConfig.Encoding = "json"

	log, err := logConfig.Build()
	if err != nil {
		log.Fatal("Failed to initialize logger")
	}

	d.log = *log

	// Загрузка конфигурации очереди
	d.consumer = consumer.NewConsumer(
		consumer.WithDB(d.db),
		consumer.WithLogger(&d.log),
	)
	RegisterJobs(d.consumer)

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
