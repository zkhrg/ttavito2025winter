package usecase

import (
	"ttavito/domain/entities"
	"ttavito/domain/interfaces"
)

type Usecase struct {
	repo interfaces.ShopRepository
}

func (u *Usecase) GetInfo(username string) (*entities.InfoResponse, error) {
	return nil, nil
}
func (u *Usecase) BuyItem(username, item string) error {
	u.repo.BuyItem("john_doe", "cup")
	return nil
}
func (u *Usecase) SendCoin(senderUsername string, recipientUsername string, amount int) error {
	u.repo.SendCoin(senderUsername, recipientUsername, amount)
	return nil
}

func NewUsecase(repo interfaces.ShopRepository) *Usecase {
	return &Usecase{repo: repo}
}
