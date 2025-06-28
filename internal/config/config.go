package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string        `mapstructure:"SERVER_PORT"`
	DBHost          string        `mapstructure:"DB_HOST"`
	DBPort          int           `mapstructure:"DB_PORT"`
	DBUser          string        `mapstructure:"DB_USER"`
	DBPassword      string        `mapstructure:"DB_PASSWORD"`
	DBName          string        `mapstructure:"DB_NAME"`
	RedisURL        string        `mapstructure:"REDIS_URL"`
	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	AccessTokenExp  time.Duration `mapstructure:"ACCESS_TOKEN_EXP"`
	RefreshTokenExp time.Duration `mapstructure:"REFRESH_TOKEN_EXP"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("ACCESS_TOKEN_EXP", 15*time.Minute)
	viper.SetDefault("REFRESH_TOKEN_EXP", 168*time.Hour) // 7 days

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
