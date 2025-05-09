package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig    `yaml:"app" env-prefix:"APP_"`
	Database DBConfig     `yaml:"db" env-prefix:"DB_"`
	Logger   LoggerConfig `yaml:"logger" env-prefix:"LOG_"`
	JWT      JWTConfig    `yaml:"jwt"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"NAME" env-default:"Client Service"`
	Port string `yaml:"port" env:"PORT" env-default:"50051"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     uint16 `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD"`
	Name     string `yaml:"name" env:"NAME"`
	Schema   string `yaml:"schema" env:"SCHEMA" env-default:"public"`
	SSLMode  string `yaml:"sslmode" env:"SSLMODE" env-default:"disable"`
	MaxConn  int32  `yaml:"max_connections" env:"MAX_CONN" env-default:"5"`
	MinConn  int32  `yaml:"min_connections" env:"MIN_CONN" env-default:"1"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LEVEL" env-default:"prod"`
}

type JWTConfig struct {
	AccessTokenDuration  time.Duration `yaml:"access_token_duration" env:"ACCESS_TOKEN_DURATION" env-default:"24h"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration" env:"TOKEN_DURATION" env-default:"168h"`
	JWTSecret            string        `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"secret"`
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
