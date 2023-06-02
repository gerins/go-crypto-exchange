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
	Database      struct {
		Read  Database
		Write Database
	}
}

type Cache struct {
	Address  string
	Password string
	Database int
}

type MessageBroker struct {
	Brokers string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Pass     string
	Name     string
	DebugLog bool
	Pool     struct {
		MaxIdleConn     int
		MaxOpenConn     int
		MaxConnLifetime time.Duration
		MaxIdleLifeTime time.Duration
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
