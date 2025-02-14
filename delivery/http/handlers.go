package http

import (
	"encoding/json"
	"net/http"

	"ttavito/domain/entities"
	"ttavito/internal"
	"ttavito/usecase"
)

func GetInfoHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func SendCoinHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(internal.ValidSendCoinKey).(entities.SendCoinRequest)
		if !ok {
			http.Error(w, "Invalid request", http.StatusInternalServerError)
			return
		}
		uc.SendCoin("zkhrg", req.ToUser, req.Amount)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("lol")
	}
}

func BuyItemHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uc.BuyItem("john_doe", "cup")
	}
}

func AuthHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
