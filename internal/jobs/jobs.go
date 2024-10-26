package jobs

import (
	"app_burse_backend/pkg/queue/consumer"
	"app_burse_backend/pkg/queue/job"

	"go.uber.org/zap"
)

func RegisterJobs(consumer *consumer.Consumer) {
	consumer.RegisterHandler("send_email", SendEmail)
}

func SendEmail(job *job.ProcessJob) error {
	// Отправка email
	params, err := job.Payload()
	if err != nil {
		job.Log().Info("Ошибка получения параметров email", zap.Error(err))
	}

	email, ok := (*params)["email"].(string)
	if !ok {
		job.Log().Info("Ошибка получения темы email", zap.Error(err))
		return nil
	}

	job.Log().Info("jjj", zap.String("email", email))
	return nil
}
