package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	configFileTypeYaml = "yaml"
)

type Config struct {
	App          App
	GRPC         GRPC
	Dependencies Dependencies
}

type App struct {
	Name       string
	Version    string
	Host       string
	Port       string
	CtxTimeout time.Duration
}

type GRPC struct {
	Host string
	Port string
}

type Dependencies struct {
	Cache         Cache
	MessageBroker MessageBroker
}

type Cache struct {
	Address  string
	Password string
	Database int
}

type MessageBroker struct {
	Brokers  string
	Group    string
	Consumer struct {
		Topic string
	}
	Producer struct {
		Topic string
	}
}

// ParseConfigFile is used for parsing config file into struct
func ParseConfigFile(configName string) *Config {
	viperConfig := viper.New()
	viperConfig.SetConfigName(configName)
	viperConfig.SetConfigType(configFileTypeYaml)
	viperConfig.AddConfigPath(".")

	if err := viperConfig.ReadInConfig(); err != nil {
		log.Fatalf("failed reading config, %v", err)
	}

	config := new(Config)
	if err := viperConfig.Unmarshal(&config); err != nil {
		log.Fatalf("failed parsing config, %v", err)
	}

	return config
}
