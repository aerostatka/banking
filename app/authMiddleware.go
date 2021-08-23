package app

import (
	"github.com/aerostatka/banking/domain"
	"github.com/aerostatka/banking-lib/errs"
	"github.com/aerostatka/banking-lib/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	repo domain.AuthRepository
}

func (am AuthMiddleware) authorizationHandler() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			currentRoute := mux.CurrentRoute(r)
			currentRouteVars := mux.Vars(r)
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				token := GetTokenFromHeader(authHeader)

				if am.repo.IsAuthorized(token, currentRoute.GetName(), currentRouteVars) {
					next.ServeHTTP(rw, r)
				} else {
					logger.Error("Token is not valid for route " + currentRoute.GetName())
					appErr := errs.NewForbiddenError("Token is not valid")
					writeResponse(rw, appErr.Code, "application/json", appErr.Message)
				}
			} else {
				logger.Error("Authorization Header is missing for route " + currentRoute.GetName())
				appErr := errs.NewUnauthorizedError("Token is missing")
				writeResponse(rw, appErr.Code, "application/json", appErr.Message)
			}
		})
	}
}

func GetTokenFromHeader(header string) string  {
	parts := strings.Split(header, " ")

	if len(parts) != 2 {
		return ""
	}

	if parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}