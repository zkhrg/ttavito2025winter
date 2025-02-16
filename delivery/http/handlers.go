package http

import (
	"encoding/json"
	"net/http"

	"ttavito/domain/entities"
	"ttavito/internal"
)

func GetInfoHandler(uc UsecaseShop) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(internal.UsernameContextKey).(string)
		if !ok {
			http.Error(w, "Can't grab username from JWT", http.StatusInternalServerError)
			return
		}
		res, err := uc.GetInfo(r.Context(), username)
		if err != nil {
			http.Error(w, "Can't get info", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(*res)
	}
}

func SendCoinHandler(uc UsecaseShop) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(internal.ValidSendCoinKey).(entities.SendCoinRequest)
		if !ok {
			http.Error(w, "Invalid request", http.StatusInternalServerError)
			return
		}

		username, ok := r.Context().Value(internal.UsernameContextKey).(string)
		if !ok {
			http.Error(w, "Can't grab username from JWT", http.StatusInternalServerError)
			return
		}
		err := uc.SendCoin(r.Context(), username, req.ToUser, req.Amount)
		if err != nil {
			http.Error(w, "Can't send coins", http.StatusInternalServerError)
			return
		}
	}
}

func BuyItemHandler(uc UsecaseShop) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, ok := r.Context().Value(internal.ValidBuyItemKey).(string)
		if !ok {
			http.Error(w, "Invalid request", http.StatusInternalServerError)
			return
		}

		username, ok := r.Context().Value(internal.UsernameContextKey).(string)
		if !ok {
			http.Error(w, "Can't grab username from JWT", http.StatusInternalServerError)
			return
		}

		err := uc.BuyItem(r.Context(), username, item)
		if err != nil {
			http.Error(w, "Can't buy item", http.StatusInternalServerError)
			return
		}
	}
}

func AuthHandler(uc UsecaseShop) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(internal.ValidAuthReqKey).(entities.AuthRequest)
		if !ok {
			http.Error(w, "Invalid request", http.StatusInternalServerError)
			return
		}
		err := uc.Auth(r.Context(), req.Username, req.Password)

		if err != nil {
			http.Error(w, "Could not generate token", http.StatusUnauthorized)
			return
		}
		jwttool := internal.JWTTool{}
		token, err := jwttool.GenerateToken(req.Username, req.Password)

		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
