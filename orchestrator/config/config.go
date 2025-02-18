package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	App    AppConfig    `yaml:"app"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type AppConfig struct {
	Name string `yaml:"name"`
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
