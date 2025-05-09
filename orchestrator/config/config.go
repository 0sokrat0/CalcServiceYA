package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
	Logger LoggerConfig `yaml:"logger" env-prefix:"LOG_"`
	JWT    JWTConfig    `yaml:"jwt"`
	Grpc   GrpcConfig   `yaml:"grpc"`
	Auth   AuthService  `yaml:"auth"`
}

type AuthService struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type GrpcConfig struct {
	Port string `yaml:"port"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LEVEL" env-default:"prod"`
}

type AppConfig struct {
	Name                   string `yaml:"name"`
	TIME_ADDITION_MS       int64  `yaml:"time_addition_ms"`
	TIME_SUBTRACTION_MS    int64  `yaml:"time_subtraction_ms"`
	TIME_MULTIPLICATION_MS int64  `yaml:"time_multiplication_ms"`
	TIME_DIVISION_MS       int64  `yaml:"time_division_ms"`
}
type JWTConfig struct {
	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"secret"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "config/config.yaml"
		}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			log.Printf("❌ Config error: %v\n%s", err, help)
			os.Exit(1)
		}
	})
	return cfg
}
