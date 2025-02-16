package repository

import (
	"context"
	"fmt"
	"log/slog"

	"ttavito/domain/entities"
	"ttavito/domain/interfaces"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type EntityRepo struct {
	db      interfaces.DB
	builder sq.StatementBuilderType
}

func NewEntityRepo(db interfaces.DB) interfaces.ShopRepository {
	slog.Info("Repository created")
	return &EntityRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *EntityRepo) BuyItem(ctx context.Context, username, item string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to start transaction", "error", err)
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			slog.Error("Failed to finish buy item transaction", "error", err)
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err == nil {
				slog.Info("Buy item transaction success")
			}
		}
	}()

	selectPriceQuery, args, _ := r.builder.Select("price").
		From("products").
		Where(sq.Eq{"product_name": item}).
		ToSql()
	var price int
	err = tx.QueryRow(ctx, selectPriceQuery, args...).Scan(&price)
	if err != nil {
		return fmt.Errorf("failed to fetch product price: %v", err)
	}

	selectBalanceQuery, args, _ := r.builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()
	var balance int
	err = tx.QueryRow(ctx, selectBalanceQuery, args...).Scan(&balance)
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
	_, err = tx.Exec(ctx, updateBalanceQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	insertPurchaseQuery, args, _ := r.builder.Insert("purchases").
		Columns("username", "product_name").
		Values(username, item).
		ToSql()
	_, err = tx.Exec(ctx, insertPurchaseQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert purchase record: %v", err)
	}

	return nil
}

func (r *EntityRepo) GetInfo(ctx context.Context, username string) (*entities.InfoResponse, error) {
	var res entities.InfoResponse

	// Получаем баланс пользователя
	q, args, _ := r.builder.
		Select("balance").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()

	err := r.db.QueryRow(ctx, q, args...).Scan(&res.Coins)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	// Получаем инвентарь пользователя
	res.Inventory, err = r.GetUserInventory(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user inventory: %w", err)
	}

	// Получаем историю транзакций
	res.CoinHistory.Sent, res.CoinHistory.Received, err = r.GetUserTransactions(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return &res, nil
}

func (r *EntityRepo) Auth(ctx context.Context, username, password string) (bool, error) {
	var existingPassword string
	q, args, _ := r.builder.Select("user_password").
		From("users").
		Where(sq.Eq{"username": username}).
		ToSql()

	err := r.db.QueryRow(ctx, q, args...).Scan(&existingPassword)

	defer func() {
		if err != nil {
			slog.Error("Failed to finish auth", "error", err)
		} else {
			slog.Info("Auth success")
		}
	}()

	if err == nil {
		if existingPassword != password {
			return false, nil
		}
		return true, nil
	}

	if err.Error() != pgx.ErrNoRows.Error() {
		return false, err
	}

	q, args, _ = r.builder.Insert("users").
		Columns("username", "user_password").
		Values(username, password).
		ToSql()

	_, err = r.db.Exec(ctx, q, args...)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *EntityRepo) SendCoin(ctx context.Context, senderUsername string, recipientUsername string, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			slog.Error("Failed to send coin", "error", err)
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err == nil {
				slog.Info("success send coin")
			}
		}
	}()

	query, args, _ := r.builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": senderUsername}).
		ToSql()

	var senderBalance int
	err = tx.QueryRow(ctx, query, args...).Scan(&senderBalance)
	if err != nil {
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
	_, err = tx.Exec(ctx, updateSenderBalance, args...)

	if err != nil {
		return fmt.Errorf("failed to update sender's balance: %v", err)
	}

	query, args, _ = r.builder.Select("balance").
		From("users").
		Where(sq.Eq{"username": recipientUsername}).
		ToSql()

	var receiverBalance int
	err = tx.QueryRow(ctx, query, args...).Scan(&receiverBalance)
	if err != nil {
		return fmt.Errorf("unable to get receiver's balance: %v", err)
	}

	updateReceiverBalance, args, _ := r.builder.Update("users").
		Set("balance", receiverBalance+amount).
		Where(sq.Eq{"username": recipientUsername}).
		ToSql()
	_, err = tx.Exec(ctx, updateReceiverBalance, args...)
	if err != nil {
		return fmt.Errorf("failed to update receiver's balance: %v", err)
	}

	transferQuery, args, _ := r.builder.Insert("transfers").
		Columns("sender_username", "receiver_username", "amount").
		Values(senderUsername, recipientUsername, amount).
		ToSql()
	_, err = tx.Exec(ctx, transferQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to add in transfers: %v", err)
	}

	return nil
}

func (r *EntityRepo) GetUserInventory(ctx context.Context, username string) ([]entities.ItemResponse, error) {
	q, args, _ := r.builder.
		Select("p.product_name", "COUNT(pu.product_name) as quantity"). // Считаем количество
		From("purchases pu").
		Join("products p ON pu.product_name = p.product_name").
		Where(sq.Eq{"pu.username": username}).
		GroupBy("p.product_name"). // Группируем по названию предмета
		ToSql()

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	defer rows.Close()

	var inventory []entities.ItemResponse
	for rows.Next() {
		var item entities.ItemResponse
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan inventory row: %w", err)
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

func (r *EntityRepo) GetUserTransactions(ctx context.Context, username string) ([]entities.SentResponse, []entities.ReceivedResponse, error) {
	// Транзакции от пользователя
	qSent, args, _ := r.builder.
		Select("receiver_username", "amount").
		From("transfers").
		Where(sq.Eq{"sender_username": username}).
		ToSql()

	rowsSent, err := r.db.Query(ctx, qSent, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rowsSent.Close()

	var sent []entities.SentResponse
	for rowsSent.Next() {
		var t entities.SentResponse
		if err := rowsSent.Scan(&t.ToUser, &t.Amount); err != nil {
			return nil, nil, err
		}
		sent = append(sent, t)
	}

	qReceived, args, _ := r.builder.
		Select("sender_username", "amount").
		From("transfers").
		Where(sq.Eq{"receiver_username": username}).
		ToSql()

	rowsReceived, err := r.db.Query(ctx, qReceived, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rowsReceived.Close()

	var received []entities.ReceivedResponse
	for rowsReceived.Next() {
		var t entities.ReceivedResponse
		if err := rowsReceived.Scan(&t.FromUser, &t.Amount); err != nil {
			return nil, nil, err
		}
		received = append(received, t)
	}

	return sent, received, nil
}
