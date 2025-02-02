package main

import (
	"log"

	"go.uber.org/zap"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/server"
	"cyansnbrst/products-service/pkg/db/postgres"
	"cyansnbrst/products-service/pkg/kafka"
)

//	@title			Products Service API
//	@version		1.0
//	@description	API Server for view and edit products

//	@host		localhost:8080
//	@BasePath	/products
//	@host

// @securityDefinitions.apikey	cookieAuth
// @in							cookie
// @name						token
func main() {
	log.Println("starting products server")

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

	kafkaUserWriter, err := kafka.InitKafkaWriter(cfg, "user")
	if err != nil {
		logger.Fatal("failed to init kafka user producer",
			zap.String("error", err.Error()),
		)
	}
	logger.Info("kafka user producer connected")

	kafkaProductWriter, err := kafka.InitKafkaWriter(cfg, "product")
	if err != nil {
		logger.Fatal("failed to init kafka product producer",
			zap.String("error", err.Error()),
		)
	}
	logger.Info("kafka product producer connected")

	s := server.NewServer(cfg, logger, psqlDB, kafkaUserWriter, kafkaProductWriter)
	if err = s.Run(); err != nil {
		logger.Fatal("an error occured",
			zap.String("error", err.Error()),
		)
	}
}
