package web

import (
	"app_burse_backend/configs"
	"app_burse_backend/pkg/logger"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/producer"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type WebContext struct {
	db        postgres.Database
	config    *configs.Config
	log       *zap.Logger
	localizer *i18n.Localizer

	producer *producer.Producer
}

func NewWebContext(config *configs.Config) *WebContext {
	return &WebContext{
		config: config,
	}
}

func (c *WebContext) InitDB() {
	c.db = postgres.NewPostgres().Connect(
		c.config.DB.Host, c.config.DB.Port, c.config.DB.Username, c.config.DB.Password, c.config.DB.Name,
	)
}

func (c *WebContext) InitLogger() {
	c.log = logger.NewLogger(logger.WithDevelopment(true), logger.WithLevel(zap.DebugLevel)).Build()
}

func (c *WebContext) InitLocales(currentPwd string) {
	// Load locales
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	bundle.MustLoadMessageFile(currentPwd + "./configs/locales/ru.yml")

	c.localizer = i18n.NewLocalizer(bundle, "ru-RU")
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

func (c *WebContext) DB() postgres.Database {
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

func (c *WebContext) Locales() *i18n.Localizer {
	return c.localizer
}
