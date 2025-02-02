package http

import (
	"net/http"

	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/middleware"
	"cyansnbrst/recommendations-service/internal/recommendations"
	erp "cyansnbrst/recommendations-service/pkg/error_responses"
	"cyansnbrst/recommendations-service/pkg/utils"
)

// Recommendations handlers
type recommendationsHandlers struct {
	cfg               *config.Config
	recommendationsUC recommendations.UseCase
	logger            *zap.Logger
}

// Recommendations handlers constructor
func NewRecommendationsHandlers(cfg *config.Config, recommendationsUC recommendations.UseCase, logger *zap.Logger) recommendations.Handlers {
	return &recommendationsHandlers{
		cfg:               cfg,
		recommendationsUC: recommendationsUC,
		logger:            logger,
	}
}

//	@Summary		Get recommendations for user
//	@Description	Retrieves personalized recommendations for the authenticated user.
//	@Tags			recommendations
//	@Produce		json
//	@Security		cookieAuth
//	@Success		200	{object}	models.RecommendationResponse	"success response with recommendations"
//	@Failure		404	{object}	models.ErrorResponse			"not found error if recommendations are unavailable"
//	@Failure		500	{object}	models.ErrorResponse			"internal server error"
//	@Router			/ [get]
func (h *recommendationsHandlers) GetInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUID := r.Context().Value(middleware.UserContextKey).(string)

		recommendations, err := h.recommendationsUC.GetRecommendationsForUser(userUID)
		if err != nil {
			erp.NotFoundResponse(w, r, h.logger)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"recommendations": recommendations,
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}
