package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sample/internal/domain/entities/config"
	authServices "sample/internal/domain/services/auth"
	"strings"
)

type UserClaimsKey struct{}

func AuthMiddlewareFactory(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Authorization header is missing"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				http.Error(w, `{"error":"Bearer token is missing"}`, http.StatusUnauthorized)
				return
			}

			token, claims, err := authServices.ValidateToken(cfg, tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}
			if token.Valid {
				log.Println("Token is valid")
			} else {
				log.Println("Token is invalid")
			}
			// ユーザー識別情報をコンテキストに追加
			key := UserClaimsKey{}
			ctx := context.WithValue(r.Context(), key, claims)
			// コンテキストを持つ新しいリクエストを作成
			r = r.WithContext(ctx)

			// Token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
