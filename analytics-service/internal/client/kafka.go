package client

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	kf "cyansnbrst/analytics-service/pkg/kafka"
)

// Kafka client struct
type KafkaClient struct {
	config  *config.Config
	logger  *zap.Logger
	readers []*readerGroup
	wg      *sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// Kafka reader group struct
type readerGroup struct {
	reader  *kafka.Reader
	handler func(kafka.Message) error
}

// New kafka client constructor
func NewKafkaClient(cfg *config.Config, logger *zap.Logger) *KafkaClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &KafkaClient{
		config:  cfg,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
		readers: make([]*readerGroup, 0),
	}
}

// Add new kafka reader
func (kc *KafkaClient) AddReader(topicKey, groupID string, handler func(kafka.Message) error) error {
	reader, err := kf.InitKafkaReader(kc.config, topicKey, groupID)
	if err != nil {
		return err
	}

	kc.readers = append(kc.readers, &readerGroup{reader: reader, handler: handler})
	return nil
}

// Run kafka client
func (kc *KafkaClient) Run() {
	kc.logger.Info("starting Kafka client")

	for _, rg := range kc.readers {
		kc.wg.Add(1)
		go func(rg *readerGroup) {
			defer kc.wg.Done()
			err := kf.ConsumeMessages(kc.ctx, rg.reader, rg.handler)
			if err != nil {
				kc.logger.Error("error consuming messages", zap.Error(err))
			}
		}(rg)
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		kc.logger.Info("shutting down Kafka client")
		kc.Stop()
	}()
}

// Stop all kafka readers
func (kc *KafkaClient) Stop() {
	kc.cancel()
	for _, rg := range kc.readers {
		if err := rg.reader.Close(); err != nil {
			kc.logger.Error("error closing Kafka reader", zap.Error(err))
		}
	}
	kc.wg.Wait()
	kc.logger.Info("Kafka client stopped")
}
