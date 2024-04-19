package server

import (
	"context"
	"github.com/golang-jwt/jwt"
	"net/http"
	"store-api/app/model"
	"store-api/app/store"
	"strings"
)

func (s *server) JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *server) UserAuthenticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			token := strings.TrimPrefix(authorization, "Token ")
			if len(token) > 0 {
				claims := jwt.MapClaims{}
				_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(s.jwtSecretKey), nil
				})
				if err == nil {
					id := claims["id"]
					if id != nil {
						u, err := s.store.User().One(model.UserOneFilter{IDData: model.IDData{ID: id.(string)}})
						if err == nil && u != nil {
							if u.Token == token {
								next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), s.jwtSecretKey, u)))
								return
							}
						}
					}
				}
			}
		}
		s.error(w, r, http.StatusUnauthorized, store.ErrUserPermissionDenied)
		return
	})
}
