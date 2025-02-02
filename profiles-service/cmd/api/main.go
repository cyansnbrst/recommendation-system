package main

import (
	"log"

	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/server"
	"cyansnbrst/profiles-service/pkg/db/postgres"
	"cyansnbrst/profiles-service/pkg/kafka"
)

//	@title			Profiles Service API
//	@version		1.0
//	@description	API Server for view and edit user's profile

//	@host		localhost:8080
//	@BasePath	/profiles
//	@host

// @securityDefinitions.apikey	cookieAuth
// @in							cookie
// @name						token
func main() {
	log.Println("starting profiles server")

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

	kafkaWriter, err := kafka.InitKafkaWriter(cfg, "user")
	if err != nil {
		logger.Fatal("failed to init kafka producer",
			zap.String("error", err.Error()),
		)
	}
	logger.Info("kafka producer connected")

	s := server.NewServer(cfg, logger, psqlDB, kafkaWriter)
	if err = s.Run(); err != nil {
		logger.Fatal("an error occured",
			zap.String("error", err.Error()),
		)
	}
}
