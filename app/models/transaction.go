package models

import (
	"apisim/app/db"
	"context"
)

type (
	Transaction struct {
		SequentialIdentifier
		UserID   int64   `json:"user_id"`
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Balance  float64 `json:"balance"`
		Code     string  `json:"code"`
		Type     string  `json:"type"`
		Timestamps
	}
)

const (
	createTransactionSQL           = `insert into transactions (user_id, amount, currency, balance, code, type, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`
	getTransactionSQL              = `select id, user_id, amount, currency, balance, code, type, created_at from transactions`
	getUserTransactionSQL          = getTransactionSQL + ` where user_id=$1`
	getLastTransactionSQL          = getUserTransactionSQL + ` order by id desc limit 1`
	getTransactionByCodeAndTypeSQL = getUserTransactionSQL + ` where code=$1 and type=$2 limit 1`
)

func (t *Transaction) ByCodeAndType(
	ctx context.Context,
	db db.SQLOperations,
	code string,
	transType string,
) (*Transaction, error) {
	row := db.QueryRowContext(ctx, getTransactionByCodeAndTypeSQL, code, transType)
	return t.scan(row)
}

func (t *Transaction) LastTransaction(
	ctx context.Context,
	db db.SQLOperations,
	userID int64,
) (*Transaction, error) {
	row := db.QueryRowContext(ctx, getLastTransactionSQL, userID)
	return t.scan(row)
}

func (r *Transaction) scan(
	row db.RowScanner,
) (*Transaction, error) {
	var transaction Transaction

	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Balance,
		&transaction.Code,
		&transaction.Type,
		&transaction.CreatedAt,
	)

	return &transaction, err
}

func (transaction *Transaction) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	transaction.Timestamps.Touch()

	err := db.QueryRowContext(
		ctx,
		createTransactionSQL,
		transaction.UserID,
		transaction.Amount,
		transaction.Currency,
		transaction.Balance,
		transaction.Code,
		transaction.Type,
		transaction.Timestamps.CreatedAt,
	).Scan(&transaction.ID)
	return err
}
