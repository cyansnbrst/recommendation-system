package main

import (
	"log"

	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/server"
	"cyansnbrst/recommendations-service/pkg/db/postgres"
	"cyansnbrst/recommendations-service/pkg/db/redis"
)

//	@title			Recommendations Service API
//	@version		1.0
//	@description	API Server for get user's recommendations

//	@host		localhost:8080
//	@BasePath	/recommendations
//	@host

// @securityDefinitions.apikey	cookieAuth
// @in							cookie
// @name						token
func main() {
	log.Println("starting recommendations server")

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

	redisClient := redis.NewRedisClient(cfg)
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("failed to close redis", zap.String("error", err.Error()))
		}
	}()
	logger.Info("redis connected")

	s := server.NewServer(cfg, logger, psqlDB, redisClient)
	if err = s.Run(); err != nil {
		logger.Fatal("an error occured",
			zap.String("error", err.Error()),
		)
	}
}
