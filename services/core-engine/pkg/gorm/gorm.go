package gorm

import (
	"fmt"

	"github.com/gerins/log"
	gormLogger "github.com/gerins/log/extension/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"core-engine/config"
)

func Init(cfg config.Database) *gorm.DB {
	address := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Pass, cfg.Name, cfg.Port, "disable")

	logMode := gormLogger.Default
	if cfg.DebugLog {
		logMode = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(address), &gorm.Config{Logger: logMode})
	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err.Error())
	}

	sqlDB.SetMaxIdleConns(cfg.Pool.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.Pool.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(cfg.Pool.MaxIdleLifeTime)
	sqlDB.SetConnMaxLifetime(cfg.Pool.MaxConnLifetime)

	log.Info("GormDB : Successfully Connected to Database")
	return db
}
