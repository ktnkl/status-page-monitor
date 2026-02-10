package middleware

import (
	"context"
	"log"
	"net/http"
	jwtc "status-page-monitor/server/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func Auth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Authorization token is not provided")
				http.Error(w, "Authorization token is not provided", http.StatusUnauthorized)
				return
			}

			if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				log.Println("Too short or doesnt have 'Bearer'")
				http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[7:]

			token, err := jwtc.ValidateJWT(tokenString)

			if err != nil || !token.Valid {
				log.Println("Error while token validation:", err.Error())
				http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := uint(claims["user_id"].(float64))
				r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))
				r = r.WithContext(context.WithValue(r.Context(), "login", claims["login"]))
			}

			f(w, r)
		}
	}
}
