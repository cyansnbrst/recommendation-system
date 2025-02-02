package erp

import (
	"net/http"

	"go.uber.org/zap"

	"cyansnbrst/auth-service/pkg/utils"
)

func logError(r *http.Request, l *zap.Logger, err error) {
	l.Error("an error occured",
		zap.String("request_method", r.Method),
		zap.String("request_url", r.URL.String()),
		zap.Error(err),
	)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}, l *zap.Logger) {
	env := utils.Envelope{"error": message}

	err := utils.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, l, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger, err error) {
	logError(r, l, err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message, l)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message, l)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), l)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	message := "invalid authentication credentials"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	message := "invalid or missing authentication token"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	message := "you must be authenticated to access this resource"
	errorResponse(w, r, http.StatusUnauthorized, message, l)
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request, l *zap.Logger) {
	message := "your user account doesnt't have the necessare permissions to access this resource"
	errorResponse(w, r, http.StatusForbidden, message, l)
}
