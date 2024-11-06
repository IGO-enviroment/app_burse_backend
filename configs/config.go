package configs

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	DB  DBConfig  `json:"db"`
	Web WebConfig `json:"web"`
}

const (
	defaultConfigsDir = "./configs"
	configType        = "yml"
	configName        = "config"
)

func NewCondfig() *Config {
	return &Config{}
}

func (c *Config) Load() *Config {
	viper.AddConfigPath(defaultConfigsDir)

	err := c.readConfigs()
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) LoadForTest(currentPwd string) *Config {
	viper.AddConfigPath(currentPwd + defaultConfigsDir)

	err := c.readConfigs()
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Config) readConfigs() error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	return nil
}
