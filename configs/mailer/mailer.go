package configs_mailer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MailerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

func NewConfig() *MailerConfig {
	return &MailerConfig{}
}

func (c *MailerConfig) Load() (*MailerConfig, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ymlFile, err := os.ReadFile(fmt.Sprintf("%s/configs/mailer/settings.yml", pwd))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(ymlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
