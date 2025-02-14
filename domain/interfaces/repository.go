package interfaces

import "ttavito/domain/entities"

type ShopRepository interface {
	GetInfo(username string) (*entities.InfoResponse, error)
	BuyItem(username, item string) error
	SendCoin(senderUsername string, recipientUsername string, amount int) error
}
