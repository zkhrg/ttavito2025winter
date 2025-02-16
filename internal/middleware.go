package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"ttavito/domain/entities"
)

type ContextKey string

type TokenValidator interface {
	ValidateToken(token string) (string, error)
}

const (
	UsernameContextKey ContextKey = "username"
	ValidSendCoinKey   ContextKey = "validSendCoinReq"
	ValidAuthReqKey    ContextKey = "validAuthReq"
	ValidBuyItemKey    ContextKey = "validBuyItemReq"
)

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Проходим по всем миддлварям в обратном порядке, чтобы
	// первый миддлварь был самым внешним, а последний - самым внутренним
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func GetMethodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PostMethodMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		jwttool := JWTTool{}
		username, err := jwttool.ValidateToken(token)
		if err != nil || username == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateSendCoinMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req entities.SendCoinRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.ToUser == "" || req.Amount <= 0 {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ValidSendCoinKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateBuyItemMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		item := r.PathValue("item")

		if item == "" {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ValidBuyItemKey, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValdateAuthRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req entities.AuthRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Username == "" || req.Password == "" {
			http.Error(w, "Invalid input data", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ValidAuthReqKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
