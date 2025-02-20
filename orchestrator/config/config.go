package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Name                   string `yaml:"name"`
	TIME_ADDITION_MS       int64  `yaml:"time_addition_ms"`
	TIME_SUBTRACTION_MS    int64  `yaml:"time_subtraction_ms"`
	TIME_MULTIPLICATION_MS int64  `yaml:"time_multiplication_ms"`
	TIME_DIVISION_MS       int64  `yaml:"time_division_ms"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
