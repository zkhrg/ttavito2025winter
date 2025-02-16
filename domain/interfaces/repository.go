package interfaces

import (
	"context"
	"ttavito/domain/entities"
)

type ShopRepository interface {
	GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error)
	BuyItem(ctx context.Context, username, item string) error
	SendCoin(ctx context.Context, senderUsername string, recipientUsername string, amount int) error
	Auth(ctx context.Context, username, password string) (bool, error)
}
