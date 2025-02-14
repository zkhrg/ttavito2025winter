package repository

import (
	"database/sql"
	"fmt"

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

func (r *EntityRepo) BuyItem(username, item string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	fmt.Println("here")

	defer func() {
		if err != nil {
			fmt.Println("buy item rollback", err)
			tx.Rollback()
		} else {
			fmt.Println("buy item commit")
			err = tx.Commit()
		}
	}()

	selectPriceQuery, args, _ := r.Builder.Select("price").
		From("products").
		Where(sq.Eq{"product_name": item}).
		ToSql()
	var price int
	err = tx.QueryRow(selectPriceQuery, args...).Scan(&price)
	if err != nil {
		return fmt.Errorf("failed to fetch product price: %v", err)
	}

	selectBalanceQuery, args, _ := r.Builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()
	var balance int
	err = tx.QueryRow(selectBalanceQuery, args...).Scan(&balance)
	if err != nil {
		return fmt.Errorf("failed to fetch user balance: %v", err)
	}

	if balance < price {
		return fmt.Errorf("not enough balance to buy the product")
	}

	updateBalanceQuery, args, _ := r.Builder.Update("users").
		Set("balance", sq.Expr("balance - ?", price)).
		Where(sq.Eq{"username": username}).
		ToSql()
	_, err = tx.Exec(updateBalanceQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	insertPurchaseQuery, args, _ := r.Builder.Insert("purchases").
		Columns("username", "product_name").
		Values(username, item).
		ToSql()
	_, err = tx.Exec(insertPurchaseQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert purchase record: %v", err)
	}

	return nil
}

func (r *EntityRepo) GetInfo(username string) (*entities.InfoResponse, error) {
	return nil, nil
}

func (r *EntityRepo) SendCoin(senderUsername, recipientUsername string) error {
	return nil
}
