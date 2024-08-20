package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	configFileTypeYaml = "yaml"
)

type Config struct {
	App          App
	Security     Security
	Dependencies Dependencies
}

type App struct {
	Name    string
	Version string
	HTTP    HTTP
	GRPC    GRPC
	Logging Logging
}

type HTTP struct {
	Host       string
	Port       string
	CtxTimeout time.Duration
}

type GRPC struct {
	Host string
	Port string
}

type Logging struct {
	LogToTerminal     bool   // Set log output to stdout
	LogToFile         bool   // Set log output to file
	Location          string // Location log file will be save in project directory
	FileLogName       string // File log name.
	MaxAge            int    // Days before deleting log file.
	RotationFile      int    // Hour before creating new file.
	HideSensitiveData bool   // Enable hide sensitive data with struct tag `log:"hide"`
}

type Security struct {
	Jwt struct {
		Key      string
		Duration time.Duration
	}
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
	Brokers  string
	Group    string
	Consumer struct {
		Topic struct {
			MatchOrder string
		}
	}
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

	// Override with environment variable if exist
	if env := os.Getenv("DATABASE_MASTER_HOST"); env != "" {
		config.Dependencies.Database.Write.Host = env
	}
	if env := os.Getenv("DATABASE_SLAVE_HOST"); env != "" {
		config.Dependencies.Database.Read.Host = env
	}
	if env := os.Getenv("REDIS_ADDRESS"); env != "" {
		config.Dependencies.Cache.Address = env
	}
	if env := os.Getenv("KAFKA_ADDRESS"); env != "" {
		config.Dependencies.MessageBroker.Brokers = env
	}

	return config
}
