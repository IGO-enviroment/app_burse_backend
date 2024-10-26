package configs

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBHost string `yaml:"db_host"`
	DBPort int    `yaml:"db_port"`

	Web WebConfig `yaml:"web"`
}

func NewCondfig() *Config {
	return &Config{}
}

func (c *Config) Load() *Config {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	ymlFile, err := os.ReadFile(fmt.Sprintf("%s/configs/config.yml", pwd))
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(ymlFile, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
