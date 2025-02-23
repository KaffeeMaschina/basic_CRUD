package apiserver

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	BindAddr string `env:"BIND_ADDR" envDefault:":8080"`
	DBName   string `env:"PG_DATABASE_NAME,required"`
	DBUser   string `env:"PG_USER,required"`
	DBPass   string `env:"PG_PASSWORD,required"`
	DBHost   string `env:"PG_HOST" envDefault:"localhost"`
	DBPort   string `env:"PG_PORT,required"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var cfg Config
	if err = env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing config: ", err)
	}
	return &cfg

}
