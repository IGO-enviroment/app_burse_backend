package web

import (
	"app_burse_backend/configs"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/producer"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type WebContext struct {
	db     *sqlx.DB
	config *configs.Config
	log    *zap.Logger

	producer *producer.Producer
}

func NewWebContext(config *configs.Config) *WebContext {
	return &WebContext{
		config: config,
	}
}

func (c *WebContext) InitDB() {
	fmt.Println("sfdf")
	c.db = postgres.NewPostgres().Connect(c.config.DBHost, c.config.DBPort, "user", "pass", "postgres")

	c.db.MustExec(`
	DROP TABLE IF EXISTS queue_jobs;

	CREATE TABLE IF NOT EXISTS queue_jobs (
    id SERIAL PRIMARY KEY,

		reserv_at TIMESTAMP,

		processed BOOLEAN DEFAULT false,
    queue_name VARCHAR(255) NOT NULL,
    run_at TIMESTAMP,
	
		method VARCHAR(255) NOT NULL,
		item TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
	);`)
}

func (c *WebContext) InitLogger() error {
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

	c.log = log
	return nil
}

func (c *WebContext) InitProducer() error {
	p := producer.NewProducerQueue(
		producer.WithLogger(c.log),
		producer.WithQueues("default"),
		producer.WithTableName("queue_jobs"),
	)

	c.producer = p
	return nil
}

func (c *WebContext) SetupRoutes() *mux.Router {
	r := NewRoutes(c)
	return r.Setup()
}

func (c *WebContext) Run() {

	r := c.SetupRoutes()

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%d", c.config.Web.Host, c.config.Web.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (c *WebContext) DB() *sqlx.DB {
	return c.db
}

func (c *WebContext) Producer() *producer.Producer {
	return c.producer
}

func (c *WebContext) Logger() *zap.Logger {
	return c.log
}

func (c *WebContext) Configs() *configs.Config {
	return c.config
}
