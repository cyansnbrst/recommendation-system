package usecase

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/models"
	"cyansnbrst/profiles-service/internal/profiles"
	kf "cyansnbrst/profiles-service/pkg/kafka"
)

// Profiles usecase struct
type profilesUC struct {
	cfg          *config.Config
	profilesRepo profiles.Repository
	logger       *zap.Logger
}

// Prodiles usecase constructor
func NewProfilesUseCase(cfg *config.Config, profilesRepo profiles.Repository, logger *zap.Logger) profiles.UseCase {
	return &profilesUC{cfg: cfg, profilesRepo: profilesRepo, logger: logger}
}

// Get a profile by UID
func (u *profilesUC) Get(uid string) (*models.Profile, error) {
	return u.profilesRepo.Get(uid)
}

// Update profile data
func (u *profilesUC) Update(uid string, location *string, interests []string) error {
	profile, err := u.profilesRepo.Get(uid)
	if err != nil {
		return err
	}

	if location != nil {
		profile.Location = *location
	}
	if interests != nil {
		profile.Interests = interests
	}

	return u.profilesRepo.Update(profile)
}

// Create profile
func (u *profilesUC) CreateProfile(uid string, name string) error {
	defaultLocation := u.cfg.DefaultLocation
	defaultInterests := []string{u.cfg.DefaultInterests}
	return u.profilesRepo.CreateProfile(uid, name, defaultLocation, defaultInterests)
}

// Send message to kafka
func (u *profilesUC) SendToKafka(ctx context.Context, key string, message kf.KafkaMessage, writer *kafka.Writer) error {
	messageValue, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMessage := kafka.Message{
		Key:   []byte(key),
		Value: messageValue,
	}

	if err = writer.WriteMessages(ctx, kafkaMessage); err != nil {
		return err
	}

	return nil
}
