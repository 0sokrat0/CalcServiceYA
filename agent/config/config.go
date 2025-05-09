package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
	Logger LoggerConfig `yaml:"logger" env-prefix:"LOG_"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LEVEL" env-default:"prod"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	Name string `yaml:"name"`

	COMPUTING_POWER int `yaml:"computing_power"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	v.AutomaticEnv()

	v.AutomaticEnv()
	v.SetEnvPrefix("AGENT")
	v.BindEnv("server.host")
	v.BindEnv("server.port")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
