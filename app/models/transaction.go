package models

import (
	"apisim/app/db"
	"apisim/app/helpers"
	"context"
	"fmt"
	"strings"
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
	createTransactionSQL           = `insert into transactions (user_id, amount, currency, balance, code, type, created_at) values ($1, $2, $3, $4, $5, $6, $7) returning id`
	getTransactionSQL              = `select id, user_id, amount, currency, balance, code, type, created_at from transactions`
	getUserTransactionSQL          = getTransactionSQL + ` where user_id=$1`
	getLastTransactionSQL          = getUserTransactionSQL + ` order by id desc limit 1`
	getTransactionByCodeAndTypeSQL = getUserTransactionSQL + ` where code=$1 and type=$2 limit 1`
	countTransactionSQL            = `select count(id) from transactions`
)

func (t *Transaction) AllForUser(
	ctx context.Context,
	db db.SQLOperations,
	userID int64,
	filter *Filter,
) ([]*Transaction, error) {
	trans := make([]*Transaction, 0)

	query, args := t.buildQuery(
		getTransactionSQL,
		userID,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return trans, err
	}

	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Balance,
			&transaction.Code,
			&transaction.Type,
			&transaction.CreatedAt,
		)
		if err != nil {
			return trans, err
		}
		trans = append(trans, &transaction)
	}

	return trans, err
}

func (t *Transaction) ByCodeAndType(
	ctx context.Context,
	db db.SQLOperations,
	code string,
	transType string,
) (*Transaction, error) {
	row := db.QueryRowContext(ctx, getTransactionByCodeAndTypeSQL, code, transType)
	return t.scan(row)
}

func (t *Transaction) Count(
	ctx context.Context,
	db db.SQLOperations,
	userID int64,
	filter *Filter,
) (int, error) {
	query, args := t.buildQuery(
		countTransactionSQL,
		userID,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
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

func (t *Transaction) buildQuery(
	query string,
	userID int64,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	conditions = append(conditions, fmt.Sprintf(" user_id = $%d", placeholder.Touch()))
	args = append(args, userID)

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}
