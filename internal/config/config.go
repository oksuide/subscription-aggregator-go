package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env              string           `mapstructure:"env"`
	HTTPServerConfig HTTPServerConfig `mapstructure:"http_server"`
	DBConfig         DBConfig         `mapstructure:"db"`
}

type HTTPServerConfig struct {
	Address     string        `mapstructure:"address"`
	ServerPort  string        `mapstructure:"server_port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"dbname"`
	SSLMode        string `mapstructure:"sslmode"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

func LoadConfig() (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
