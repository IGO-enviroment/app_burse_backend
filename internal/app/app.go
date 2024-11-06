package app

import (
	"app_burse_backend/configs"
	"app_burse_backend/pkg/postgres"
	"app_burse_backend/pkg/queue/producer"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type AppContext interface {
	DB() postgres.Database
	Logger() *zap.Logger
	Producer() *producer.Producer
	Configs() *configs.Config
	Locales() *i18n.Localizer
}
