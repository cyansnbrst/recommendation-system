package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/middleware"
	"cyansnbrst/profiles-service/internal/models"
	"cyansnbrst/profiles-service/internal/profiles"
	"cyansnbrst/profiles-service/pkg/db"
	erp "cyansnbrst/profiles-service/pkg/error_responses"
	kf "cyansnbrst/profiles-service/pkg/kafka"
	"cyansnbrst/profiles-service/pkg/utils"
)

// Validation errors
var errNameRequired = errors.New("name is required")

// Profiles handlers
type profilesHandlers struct {
	cfg         *config.Config
	profilesUC  profiles.UseCase
	logger      *zap.Logger
	kafkaWriter *kafka.Writer
}

// Profiles handlers constructor
func NewProfilesHandlers(cfg *config.Config, profilesUC profiles.UseCase, logger *zap.Logger, kafkaWriter *kafka.Writer) profiles.Handlers {
	return &profilesHandlers{
		cfg:         cfg,
		profilesUC:  profilesUC,
		logger:      logger,
		kafkaWriter: kafkaWriter,
	}
}

// @Summary		Get user's profile info
// @Description	Retrieves user's profile info, including location and preferences.
// @Tags			profiles
// @Produce		json
// @Security		cookieAuth
// @Success		200	{object}	models.ProfileResponse	"success response with profile"
// @Failure		404	{object}	models.ErrorResponse	"not found error"
// @Failure		500	{object}	models.ErrorResponse	"internal server error"
// @Router			/ [get]
func (h *profilesHandlers) GetInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := r.Context().Value(middleware.UserContextKey).(string)

		profile, err := h.profilesUC.Get(userUID)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"profile": profile,
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

// @Summary		Edit user's profile info
// @Description	Retrieves user's profile info, including location and preferences.
// @Tags			profiles
// @Accept			json
// @Produce		json
// @Security		cookieAuth
// @Success		200	{object}	models.SuccessResponse	"success"
// @Failure		400	{object}	models.ErrorResponse	"bad request error"
// @Failure		404	{object}	models.ErrorResponse	"not found error"
// @Failure		500	{object}	models.ErrorResponse	"internal server error"
// @Router			/edit [put]
func (h *profilesHandlers) EditData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := r.Context().Value(middleware.UserContextKey).(string)

		var requestData models.EditProfileDTO

		if err := utils.ReadJSON(w, r, &requestData); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		err := h.profilesUC.Update(userUID, requestData.Location, requestData.Interests)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				erp.NotFoundResponse(w, r, h.logger)
			} else {
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		if requestData.Interests != nil {
			messagePayload := kf.KafkaMessage{
				Action: "user_update",
				Time:   time.Now().Format(time.RFC3339),
				Tags:   requestData.Interests,
			}

			err = h.profilesUC.SendToKafka(r.Context(), userUID, messagePayload, h.kafkaWriter)
			if err != nil {
				erp.ServerErrorResponse(w, r, h.logger, err)
				return
			}
		}

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"message": "profile updated successfully",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

// Create new user profile
func (h *profilesHandlers) CreateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := utils.ReadUIDParam(r)
		if userUID == "" {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		var requestData models.CreateProfileDTO

		if err := utils.ReadJSON(w, r, &requestData); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		if requestData.Name == "" {
			erp.BadRequestResponse(w, r, h.logger, errNameRequired)
			return
		}

		err := h.profilesUC.CreateProfile(userUID, requestData.Name)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}

		err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
			"message": "profile created successfully",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}
