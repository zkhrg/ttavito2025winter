package repository

import (
	"fmt"

	"database/sql"

	"ttavito/domain/entities"
	"ttavito/domain/interfaces"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type EntityRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewEntityRepo(db *sql.DB) interfaces.ShopRepository {
	return &EntityRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EntityRepo) BuyItem(username, item string) error {
	tx, err := r.db.Begin()
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

	selectPriceQuery, args, _ := r.builder.Select("price").
		From("products").
		Where(sq.Eq{"product_name": item}).
		ToSql()
	var price int
	err = tx.QueryRow(selectPriceQuery, args...).Scan(&price)
	if err != nil {
		return fmt.Errorf("failed to fetch product price: %v", err)
	}

	selectBalanceQuery, args, _ := r.builder.Select("balance").
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

	updateBalanceQuery, args, _ := r.builder.Update("users").
		Set("balance", sq.Expr("balance - ?", price)).
		Where(sq.Eq{"username": username}).
		ToSql()
	_, err = tx.Exec(updateBalanceQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	insertPurchaseQuery, args, _ := r.builder.Insert("purchases").
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

func (r *EntityRepo) Auth(username, password string) (bool, error) {
	var hashedPassword string
	q, args, _ := r.builder.Select("user_password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	err := r.db.QueryRow(q, args...).Scan(&hashedPassword)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			fmt.Println("passwords neq")
			return false, nil
		}
		return true, nil
	}

	if err != sql.ErrNoRows {
		return false, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	q, args, _ = r.builder.Insert("users").
		Columns("username", "user_password").
		Values(username, hashedPass).
		ToSql()

	_, err = r.db.Exec(q, args...)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *EntityRepo) SendCoin(senderUsername string, recipientUsername string, amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query, args, _ := r.builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": senderUsername}).
		ToSql()

	var senderBalance int
	err = tx.QueryRow(query, args...).Scan(&senderBalance)
	if err != nil {
		fmt.Println("unable to get sender's balance: sender ", senderUsername)
		return fmt.Errorf("unable to get sender's balance: %v sender %v", err, senderBalance)
	}

	if senderBalance < amount {
		err = fmt.Errorf("not enough balance for the transfer")
		return err
	}

	updateSenderBalance, args, _ := r.builder.Update("users").
		Set("balance", senderBalance-amount).
		Where(sq.Eq{"username": senderUsername}).
		ToSql()
	_, err = tx.Exec(updateSenderBalance, args...)

	if err != nil {
		return fmt.Errorf("failed to update sender's balance: %v", err)
	}

	query, args, _ = r.builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": recipientUsername}).
		ToSql()

	var receiverBalance int
	err = tx.QueryRow(query, args...).Scan(&receiverBalance)
	if err != nil {
		return fmt.Errorf("unable to get receiver's balance: %v", err)
	}

	updateReceiverBalance, args, _ := r.builder.Update("users").
		Set("balance", receiverBalance+amount).
		Where(sq.Eq{"username": recipientUsername}).
		ToSql()
	_, err = tx.Exec(updateReceiverBalance, args...)
	if err != nil {
		return fmt.Errorf("failed to update receiver's balance: %v", err)
	}

	transferQuery, args, _ := r.builder.Insert("transfers").
		Columns("sender_username", "receiver_username", "amount").
		Values(senderUsername, recipientUsername, amount).
		ToSql()
	_, err = tx.Exec(transferQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to add in transfers: %v", err)
	}

	return nil
}
