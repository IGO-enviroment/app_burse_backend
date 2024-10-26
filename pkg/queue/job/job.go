package job

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type RawJob struct {
	ID int `json:"id"`

	Method  string `json:"method"`
	Payload string `json:"payload"`

	Params map[string]interface{} `json:"-"`
}

type ProcessJob struct {
	job *RawJob

	db  *sqlx.DB
	log *zap.Logger
}

func NewProcessJob(job *RawJob, db *sqlx.DB, log *zap.Logger) *ProcessJob {
	return &ProcessJob{job: job, db: db, log: log}
}

func (p *ProcessJob) DB() *sqlx.DB {
	return p.db
}

func (p *ProcessJob) Log() *zap.Logger {
	return p.log
}

func (p *ProcessJob) Payload() (*map[string]interface{}, error) {
	err := json.Unmarshal([]byte(p.job.Payload), &p.job.Params)
	if err != nil {
		p.log.Log(zap.ErrorLevel, "Ошибка десериализации данных", zap.Error(err))
		return nil, nil
	}

	return &p.job.Params, nil
}

func (p *ProcessJob) Start(methods *map[string]func(*ProcessJob) error) error {
	// Выполнение обработки
	handler, ok := (*methods)[p.job.Method]
	if !ok {
		p.log.Log(zap.ErrorLevel, "Отсутствует обработчик для этого метода")
		return fmt.Errorf("отсутствует обработчик для этого метода: %s", p.job.Method)
	}

	err := handler(p)
	if err != nil {
		p.AfterFailure(err)
		return err
	}

	return nil
}

func (p *ProcessJob) AfterFailure(err error) {
	// Обработка ошибки
	p.log.Log(zap.ErrorLevel, "Ошибка обработки задания", zap.Error(err))
}
