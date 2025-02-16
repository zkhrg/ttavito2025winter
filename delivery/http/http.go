package http

import (
	"context"
	"net/http"

	"ttavito/domain/entities"
	"ttavito/internal"
)

type UsecaseShop interface {
	GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error)
	BuyItem(ctx context.Context, username, item string) error
	SendCoin(ctx context.Context, senderUsername string, recipientUsername string, amount int) error
	Auth(ctx context.Context, username, password string) error
}

func SetupRoutes(api UsecaseShop, mux *http.ServeMux) {
	// Создаем цепочку миддлварей и передаем API через замыкание
	buyItemCompleteHandler := internal.ChainMiddleware(
		BuyItemHandler(api),
		internal.GetMethodMiddleware,
		internal.AuthMiddleware,
		internal.ValidateBuyItemMiddleware,
	)

	authUserCompleteHandler := internal.ChainMiddleware(
		AuthHandler(api),
		internal.PostMethodMiddleware,
		internal.ValdateAuthRequestMiddleware,
	)

	sendCoinCompleteHandler := internal.ChainMiddleware(
		SendCoinHandler(api),
		internal.PostMethodMiddleware,
		internal.ValidateSendCoinMiddleware,
		internal.AuthMiddleware,
	)

	getInfoCompleteHandler := internal.ChainMiddleware(
		GetInfoHandler(api),
		internal.GetMethodMiddleware,
		internal.AuthMiddleware,
	)

	mux.Handle("/api/buy/{item}", buyItemCompleteHandler) // get
	mux.Handle("/api/auth", authUserCompleteHandler)      // post
	mux.Handle("/api/sendCoin", sendCoinCompleteHandler)  // post
	mux.Handle("/api/info", getInfoCompleteHandler)       // get
}
