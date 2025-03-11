package configuration

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"os"
)

// Configuration for client settings
type Configuration struct {
	ServerAddress string `json:"server_address"`
}

type envs struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

// Configure loads configuration from flags and env variables
func Configure() *Configuration {
	var config Configuration

	flag.StringVar(&config.ServerAddress, "a", "localhost:50051", "Server address")

	envVariables := envs{}
	err := env.Parse(&envVariables)
	if err != nil {
		zap.L().Error("Failed to parse environment variables", zap.Error(err))
	}

	_, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		config.ServerAddress = envVariables.ServerAddress
	}

	return &config
}
