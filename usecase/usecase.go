package usecase

import (
	"context"
	"fmt"
	"ttavito/domain/entities"
	"ttavito/domain/interfaces"
)

type Usecase struct {
	repo interfaces.ShopRepository
}

func (u *Usecase) GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error) {
	return u.repo.GetInfo(ctx, username)
}
func (u *Usecase) BuyItem(ctx context.Context, username, item string) error {
	return u.repo.BuyItem(ctx, username, item)
}

func (u *Usecase) SendCoin(ctx context.Context, senderUsername string, recipientUsername string, amount int) error {
	return u.repo.SendCoin(ctx, senderUsername, recipientUsername, amount)
}
func (u *Usecase) Auth(ctx context.Context, username, password string) error {
	sd, err := u.repo.Auth(ctx, username, password)
	if !sd || err != nil {
		return fmt.Errorf("passwords does not match")
	}
	return nil
}

func NewUsecase(repo interfaces.ShopRepository) *Usecase {
	return &Usecase{repo: repo}
}
