package configuration

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"os"
)

type Configuration struct {
	GRPCPort         int    `json:"grpc_port"`
	DatabaseDsn      string `json:"database_dsn"`
	JwtSecret        string `json:"jwt_secret"`
	JwtLifetimeHours int    `json:"jwt_lifetime_hours"`
}

type envs struct {
	GRPCPort         int    `env:"GRPC_PORT"`
	DatabaseDsn      string `env:"DATABASE_DSN"`
	JwtSecret        string `env:"JWT_SECRET"`
	JwtLifetimeHours int    `env:"JWT_LIFETIME_HOURS"`
}

func Configure() *Configuration {
	var config Configuration

	flag.IntVar(&config.GRPCPort, "p", 50051, "GRPC port")
	flag.StringVar(&config.DatabaseDsn, "d", "", "Адрес подключения к базе данных")
	flag.StringVar(&config.JwtSecret, "j", "secret", "JWT секрет")
	flag.IntVar(&config.JwtLifetimeHours, "l", 24, "Время жизни JWT токена в часах")
	flag.Parse()

	envVariables := envs{}
	err := env.Parse(&envVariables)
	if err != nil {
		zap.L().Error("Failed to parse environment variables", zap.Error(err))
	}

	_, exists := os.LookupEnv("GRPC_PORT")
	if exists && envVariables.GRPCPort != 0 {
		config.GRPCPort = envVariables.GRPCPort
	}

	_, exists = os.LookupEnv("DATABASE_DSN")
	if exists {
		config.DatabaseDsn = envVariables.DatabaseDsn
	}

	_, exists = os.LookupEnv("JWT_SECRET")
	if exists {
		config.JwtSecret = envVariables.JwtSecret
	}

	_, exists = os.LookupEnv("JWT_LIFETIME_HOURS")
	if exists {
		config.JwtLifetimeHours = envVariables.JwtLifetimeHours
	}

	return &config
}
