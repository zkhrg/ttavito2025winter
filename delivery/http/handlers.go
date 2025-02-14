package http

import (
	"net/http"

	"ttavito/usecase"
)

func GetInfoHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetInfo(uc, w, r)
	}
}

func SendCoinHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendCoin(uc, w, r)
	}
}

func BuyItemHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		BuyItem(uc, w, r)
	}
}

func AuthHandler(uc *usecase.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Auth(uc, w, r)
	}
}

func GetInfo(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {

}
func SendCoin(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {

}
func BuyItem(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {
	uc.BuyItem("john_doe", "cup")
}
func Auth(uc *usecase.Usecase, w http.ResponseWriter, r *http.Request) {

}
