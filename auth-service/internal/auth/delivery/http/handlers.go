package http

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"cyansnbrst/auth-service/config"
	"cyansnbrst/auth-service/internal/auth"
	"cyansnbrst/auth-service/internal/models"
	"cyansnbrst/auth-service/pkg/db"
	erp "cyansnbrst/auth-service/pkg/error_responses"
	"cyansnbrst/auth-service/pkg/utils"
)

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger *zap.Logger
}

// Auth handlers constructor
func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, logger *zap.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: logger}
}

// @Summary		Register user
// @Description	Creates a new user.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.SuccessResponse	"succesful registration"
// @Failure		400	{object}	models.ErrorResponse	"bad request error"
// @Failure		500	{object}	models.ErrorResponse	"internal server error"
// @Router			/register [post]
func (h *authHandlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.RegisterUserDTO

		if err := utils.ReadJSON(w, r, &requestBody); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		token, userUID, err := h.authUC.Create(requestBody.Email, requestBody.Password)
		fmt.Println(token, userUID, err)
		if err != nil {
			switch err {
			case db.ErrDuplicateEmail:
				erp.BadRequestResponse(w, r, h.logger, err)
			default:
				erp.ServerErrorResponse(w, r, h.logger, err)
			}
			return
		}

		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   h.cfg.Env == "production",
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(h.cfg.Timeout.Cookie),
		}

		if err = h.authUC.CreateProfile(userUID, requestBody.Name, cookie); err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		http.SetCookie(w, cookie)

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"message": "successful registration",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

// @Summary		Login user
// @Description	Generates auth token for the user.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.SuccessResponse	"successful login"
// @Failure		400	{object}	models.ErrorResponse	"bad request error"
// @Failure		401	{object}	models.ErrorResponse	"invalid authentication credentials"
// @Failure		500	{object}	models.ErrorResponse	"internal server error"
// @Router			/login [post]
func (h *authHandlers) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.LoginUserDTO

		if err := utils.ReadJSON(w, r, &requestBody); err != nil {
			erp.BadRequestResponse(w, r, h.logger, err)
			return
		}

		user, err := h.authUC.ValidateCredentials(requestBody.Email, requestBody.Password)
		if err != nil {
			erp.InvalidCredentialsResponse(w, r, h.logger)
			return
		}

		token, err := h.authUC.GenerateJWT(*user)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   h.cfg.Env == "production",
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(h.cfg.Timeout.Cookie),
		})

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"message": "successful login",
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}

// @Summary		Authenticate user
// @Description	Validates and retrieves user's information about token.
// @Tags			auth
// @Produce		json
// @Security		cookieAuth
// @Success		200	{object}	models.UserResponse		"successful login"
// @Failure		401	{object}	models.ErrorResponse	"invalid or missing authentication token"
// @Failure		500	{object}	models.ErrorResponse	"internal server error"
// @Router			/authenticate [get]
func (h *authHandlers) TokenValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			erp.InvalidAuthenticationTokenResponse(w, r, h.logger)
			return
		}

		tokenString := cookie.Value

		userID, isAdmin, err := h.authUC.ValidateToken(tokenString)
		if err != nil {
			erp.InvalidAuthenticationTokenResponse(w, r, h.logger)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{
			"user_uid": userID,
			"is_admin": isAdmin,
		}, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, h.logger, err)
		}
	}
}
