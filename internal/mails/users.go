package mails

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"go.uber.org/zap"
)

type Option func(*UserMailer)

func WithLogger(logger *zap.Logger) Option {
	return func(m *UserMailer) {
		m.log = logger
	}
}

const defaultLayout = "layouts/users.html"

type UserMailer struct {
	templatePath string
	layout       string
	from         string

	password string
	host     string
	port     int
	user     string

	log *zap.Logger
}

func NewUserMailer(opts ...Option) *UserMailer {
	m := &UserMailer{
		templatePath: "users",
		layout:       defaultLayout,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *UserMailer) SendVerificationEmail(to, subject, body string) error {
	//...
	return nil
}

func (m *UserMailer) SendResetPasswordEmail(to, subject, body string) error {
	//...
	return nil
}

func (m *UserMailer) SendPasswordChangeEmail(to, subject, body string) error {
	//...
	return nil
}

type SendWelcomeEmailData struct {
	email string
}

func (m *UserMailer) SendWelcomeEmail(to []string, data SendWelcomeEmailData) error {
	const fileName = "send_welcome_email.html"

	body, err := m.parseTemplate(fileName, data)
	if err != nil {
		return err
	}

	return m.sendEmail(to, &body)
}

// Отправка письма с использованием SMTP протокола.
func (m *UserMailer) sendEmail(to []string, body *[]byte) error {
	auth := smtp.PlainAuth("", m.user, m.password, m.host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", m.host, m.port), auth, m.from, to, *body)

	if err != nil {
		m.log.Fatal("Ошибка при отправке письма")
	}

	return nil
}

func (m *UserMailer) parseTemplate(fileName string, data any) ([]byte, error) {
	tmpl, err := template.ParseFiles(m.templatePath+"/"+fileName, m.layout)
	if err != nil {
		m.log.Fatal("Ошибка при парсинг шаблона")
		return nil, err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		m.log.Fatal("Ошибка при выполнении шаблона")
		return nil, err
	}

	return body.Bytes(), nil
}
