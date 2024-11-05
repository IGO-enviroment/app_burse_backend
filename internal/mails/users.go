package mails

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

const usersLayout = "internal/mails/layout/users.html"
const usersTemplatePath = "internal/mails/users"

type UserMailer struct {
	mailer Mailer

	log *zap.Logger
}

func NewUserMailer(log *zap.Logger) (*UserMailer, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	mailer, err := NewBaseMailer(
		WithLayout(fmt.Sprintf("%s/%s", pwd, usersLayout)),
		WithTemplatePath(fmt.Sprintf("%s/%s", pwd, usersTemplatePath)),
		WithLogger(log),
	)
	if err != nil {
		return nil, err
	}

	return &UserMailer{
		mailer: mailer,
		log:    log,
	}, nil
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
	Email    string
	Name     string
	Password string

	LinkChangePassword string
	LinkLogin          string
	LinkToProfile      string
}

// Отправка приветственного письма
func (m *UserMailer) SendWelcomeEmail(to []string, data SendWelcomeEmailData) error {
	const fileName = "send_welcome_email.html"

	body, err := m.mailer.ParseTemplate(fileName, data, true)
	if err != nil {
		return err
	}

	return m.mailer.Send(to, &body, fileName)
}
