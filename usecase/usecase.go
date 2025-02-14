package usecase

import (
	"ttavito/domain/entities"
	"ttavito/domain/interfaces"
)

type Usecase struct {
	Repo interfaces.ShopRepository
}

func GetInfo(username string) (*entities.InfoResponse, error) {
	return nil, nil
}
func BuyItem(username, item string) error {
	return nil
}
func SendCoin(senderUsername, recipientUsername string) error {
	return nil
}

func NewUsecase(repo interfaces.ShopRepository) *Usecase {
	return &Usecase{Repo: repo}
}
