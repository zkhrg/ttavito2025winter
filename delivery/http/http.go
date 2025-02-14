package http

import (
	"net/http"

	"ttavito/internal"
	"ttavito/usecase"
)

func SetupRoutes(api *usecase.Usecase, mux *http.ServeMux) {
	// Создаем цепочку миддлварей и передаем API через замыкание
	buyItemCompleteHandler := internal.ChainMiddleware(
		BuyItemHandler(api),
		internal.GetMethodMiddleware,
	)

	authUserCompleteHandler := internal.ChainMiddleware(
		AuthHandler(api),
		internal.PostMethodMiddleware,
	)

	sendCoinCompleteHandler := internal.ChainMiddleware(
		SendCoinHandler(api),
		internal.PostMethodMiddleware,
	)

	getInfoCompleteHandler := internal.ChainMiddleware(
		GetInfoHandler(api),
		internal.GetMethodMiddleware,
	)

	mux.Handle("/api/buy/{item}", buyItemCompleteHandler) // get
	mux.Handle("/api/auth", authUserCompleteHandler)      // post
	mux.Handle("/api/sendCoin", sendCoinCompleteHandler)  // post
	mux.Handle("/api/info", getInfoCompleteHandler)       // get
}
