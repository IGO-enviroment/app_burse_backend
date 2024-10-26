package app

import (
	"app_burse_backend/configs"
	"app_burse_backend/pkg/queue/producer"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AppContext interface {
	DB() *sqlx.DB
	Logger() *zap.Logger
	Producer() *producer.Producer
	Configs() *configs.Config
}
