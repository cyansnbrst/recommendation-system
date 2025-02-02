package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	erp "cyansnbrst/products-service/pkg/error_responses"
)

// Authentication middleware
func (mw *MiddlewareManager) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		req, err := http.NewRequest("GET", mw.cfg.AuthURL, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		for _, cookie := range r.Cookies() {
			req.AddCookie(cookie)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}
		defer resp.Body.Close()

		var envelope struct {
			UserUID string `json:"user_uid"`
			IsAdmin bool   `json:"is_admin"`
		}

		err = json.NewDecoder(resp.Body).Decode(&envelope)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		r = ContextSetUserUID(r, envelope.UserUID)

		next.ServeHTTP(w, r)
	})
}

// Require authentication middleware
func (mw *MiddlewareManager) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUID := ContextGetUserUID(r)
		if userUID == "" {
			erp.AuthenticationRequiredResponse(w, r, mw.logger)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Require admin rights middleware
func (mw *MiddlewareManager) RequireAdminRights(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", mw.cfg.AuthURL, nil)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		for _, cookie := range r.Cookies() {
			req.AddCookie(cookie)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				mw.logger.Error("failed to close response body", zap.Error(err))
			}
		}()

		var envelope struct {
			UserUID string `json:"user_uid"`
			IsAdmin bool   `json:"is_admin"`
		}

		err = json.NewDecoder(resp.Body).Decode(&envelope)
		if err != nil {
			erp.ServerErrorResponse(w, r, mw.logger, err)
			return
		}

		if !envelope.IsAdmin {
			erp.NotPermittedResponse(w, r, mw.logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}
