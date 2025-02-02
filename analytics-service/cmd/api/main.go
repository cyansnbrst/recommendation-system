package main

import (
	"log"

	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	"cyansnbrst/analytics-service/internal/server"
	"cyansnbrst/analytics-service/pkg/db/postgres"
)

// Run application
func main() {
	log.Println("starting analytics server")

	cfgFile, err := config.LoadConfig("config/config-local")
	if err != nil {
		log.Fatalf("loadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("parseConfig: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("failed to sync logger: %v", err)
		}
	}()

	psqlDB, err := postgres.OpenDB(cfg)
	if err != nil {
		logger.Fatal("failed to init storage",
			zap.String("error", err.Error()),
		)
	}
	defer func() {
		if err := psqlDB.Close(); err != nil {
			logger.Warn("failed to close database", zap.String("error", err.Error()))
		}
	}()
	logger.Info("database connected")

	s := server.NewServer(cfg, logger, psqlDB)
	if err = s.Run(); err != nil {
		logger.Fatal("an error occured",
			zap.String("error", err.Error()),
		)
	}
}
