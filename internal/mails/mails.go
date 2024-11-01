package mails

import (
	configs_mailer "app_burse_backend/configs/mailer"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"go.uber.org/zap"
)

type Mailer interface {
	Send(to []string, body *[]byte, typeEmail string) error
	ParseTemplate(fileName string, data any, logErrors bool) ([]byte, error)
	LoadFromConfig() error
}

type OptionBaseMailer func(*BaseMailer)

func WithTemplatePath(path string) OptionBaseMailer {
	return func(m *BaseMailer) {
		m.templatePath = path
	}
}

func WithLayout(layout string) OptionBaseMailer {
	return func(m *BaseMailer) {
		m.layout = layout
	}
}

func WithLogger(log *zap.Logger) OptionBaseMailer {
	return func(m *BaseMailer) {
		m.log = log
	}
}

type BaseMailer struct {
	templatePath string
	layout       string
	from         string

	password string
	host     string
	port     int
	user     string

	logPrefix string
	log       *zap.Logger
}

func NewBaseMailer(options ...OptionBaseMailer) (*BaseMailer, error) {
	base := &BaseMailer{
		logPrefix: "[Mailer]",
	}
	err := base.LoadFromConfig()
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(base)
	}

	return base, nil
}

// Загрузка конфигурации из файла.
func (b *BaseMailer) LoadFromConfig() error {
	cfg := configs_mailer.NewConfig()
	_, err := cfg.Load()
	if err != nil {
		return err
	}

	b.password = cfg.Password
	b.host = cfg.Host
	b.port = cfg.Port
	b.user = cfg.User
	b.from = cfg.From

	return nil
}

// Отправка письма с использованием SMTP протокола.
func (b *BaseMailer) Send(to []string, body *[]byte, typeEmail string) error {
	auth := smtp.PlainAuth("", b.user, b.password, b.host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", b.host, b.port), auth, b.from, to, *body)

	if err != nil {
		b.log.Error(
			fmt.Sprintf("%s Ошибка при отправке письма", b.logPrefix),
			zap.Error(err), zap.String("to", fmt.Sprintf("%+v", to)),
			zap.String("type", typeEmail),
		)
		return err
	}

	return nil
}

// Парсинг шаблона и подстановка переменных в шаблон.
func (b *BaseMailer) ParseTemplate(fileName string, data any, logErrors bool) ([]byte, error) {
	tmpl, err := template.ParseFiles(b.layout, b.templatePath+"/"+fileName)

	if err != nil {
		if logErrors {
			b.log.Error(
				fmt.Sprintf("%s Ошибка при парсинг шаблона", b.logPrefix),
				zap.Error(err),
				zap.String("template", fileName),
			)
		}
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		if logErrors {
			b.log.Error(
				fmt.Sprintf("%s Ошибка при выполнении шаблона", b.logPrefix),
				zap.Error(err),
				zap.String("template", fileName),
			)
		}
		return nil, err
	}

	return buf.Bytes(), nil
}
