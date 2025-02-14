package repository

import (
	"database/sql"

	"ttavito/domain/entities"
	"ttavito/domain/interfaces"

	sq "github.com/Masterminds/squirrel"
)

type EntityRepo struct {
	DB      *sql.DB
	Builder sq.StatementBuilderType
}

func NewEntityRepo(db *sql.DB) interfaces.ShopRepository {
	return &EntityRepo{
		DB:      db,
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EntityRepo) GetInfo(username string) (*entities.InfoResponse, error) {
	return nil, nil
}
func (r *EntityRepo) BuyItem(username, item string) error {
	return nil
}
func (r *EntityRepo) SendCoin(senderUsername, recipientUsername string) error {
	return nil
}
