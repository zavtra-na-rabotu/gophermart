package middleware

import (
	"context"
	"github.com/zavtra-na-rabotu/gophermart/internal/security"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	authorizationPrefix = "Bearer "
	UserIDKey           = "userId"
)

func AuthorizationMiddleware(jwtService *security.JwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeader)
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, authorizationPrefix) {
				http.Error(w, "Wrong authorization header format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, authorizationPrefix)

			claims, err := jwtService.ValidateJwtToken(token)
			if err != nil {
				zap.L().Error("Error validating token", zap.String("token", token), zap.Error(err))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Добавляем пользователя в контекст запроса
			//userID := claims.UserId
			//if err != nil {
			//	zap.L().Error("Error extracting subject", zap.String("token", token), zap.Error(err))
			//	http.Error(w, err.Error(), http.StatusUnauthorized)
			//	return
			//}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
