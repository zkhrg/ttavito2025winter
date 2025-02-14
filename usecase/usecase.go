package usecase

import (
	"ttavito/domain/entities"
	"ttavito/domain/interfaces"
)

type Usecase struct {
	Repo interfaces.ShopRepository
}

func (u *Usecase) GetInfo(username string) (*entities.InfoResponse, error) {
	return nil, nil
}
func (u *Usecase) BuyItem(username, item string) error {
	u.Repo.BuyItem("john_doe", "cup")
	return nil
}
func (u *Usecase) SendCoin(senderUsername, recipientUsername string) error {
	return nil
}

func NewUsecase(repo interfaces.ShopRepository) *Usecase {
	return &Usecase{Repo: repo}
}
