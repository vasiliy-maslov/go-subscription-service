package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LoadConfig читает конфигурацию из файла или переменных окружения.
func LoadConfig() (*PostgresConfig, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	log.Println("Ищу конфиг в папке ./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Viper не смог прочитать конфиг: %v", err)
		return nil, fmt.Errorf("ошибка чтения конфига: %w", err)
	}
	log.Println("Конфиг успешно прочитан!")

	var cfg PostgresConfig
	if err := viper.Sub("postgres").Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка десериализации конфига: %w", err)
	}

	return &cfg, nil
}
