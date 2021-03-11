package models

import (
	"apisim/app/db"
	"apisim/app/helpers"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	createMessageSQL  = `insert into messages (user_id, sender_id, meta, message, cost, currency, sent_at, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	selectMessageSQL  = `select id, user_id, sender_id, meta, message, cost, currency, sent_at, created_at from messages`
	getMessageByIDSQL = selectMessageSQL + ` where id=$1`
	countMessageSQL   = `select count(id) from messages`
)

type (
	Message struct {
		SequentialIdentifier
		UserID         int64     `json:"user_id"`
		SenderID       string    `json:"sender_id"`
		Meta           string    `json:"meta"`
		Message        string    `json:"message"`
		Cost           float64   `json:"cost"`
		Currency       string    `json:"currency"`
		SentAt         time.Time `json:"sent_at"`
		RecipientCount int       `json:"recipient_count"`
		Timestamps
	}
)

func (m *Message) AllForUser(
	ctx context.Context,
	db db.SQLOperations,
	userID int64,
	filter *Filter,
) ([]*Message, error) {
	messages := make([]*Message, 0)

	query, args := m.buildQuery(
		selectMessageSQL,
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
		return messages, err
	}

	for rows.Next() {
		var message Message
		err = rows.Scan(
			&message.ID,
			&message.UserID,
			&message.SenderID,
			&message.Meta,
			&message.Message,
			&message.Cost,
			&message.Currency,
			&message.SentAt,
			&message.CreatedAt,
		)
		if err != nil {
			return messages, err
		}
		messages = append(messages, &message)
	}

	return messages, err
}

func (m *Message) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Message, error) {
	var message Message
	row := db.QueryRowContext(ctx, getMessageByIDSQL, id)

	err := m.scan(row, &message)
	return &message, err
}

func (m *Message) Count(
	ctx context.Context,
	db db.SQLOperations,
	userID int64,
	filter *Filter,
) (int, error) {
	query, args := m.buildQuery(
		countMessageSQL,
		userID,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (m *Message) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	m.Timestamps.Touch()
	err := db.QueryRowContext(
		ctx,
		createMessageSQL,
		m.UserID,
		m.SenderID,
		m.Meta,
		m.Message,
		m.Cost,
		m.Currency,
		m.SentAt,
		m.Timestamps.CreatedAt,
	).Scan(&m.ID)
	return err
}

func (*Message) scan(
	row *sql.Row,
	message *Message,
) error {
	return row.Scan(
		&message.ID,
		&message.UserID,
		&message.SenderID,
		&message.Meta,
		&message.Message,
		&message.Cost,
		&message.Currency,
		&message.SentAt,
		&message.CreatedAt,
	)
}

func (m *Message) buildQuery(
	query string,
	userID int64,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"message"}
		for _, col := range columns {
			search := fmt.Sprintf(" (lower(%s) like '%%' || $%d || '%%')", col, placeholder.Touch())
			likeStmt = append(likeStmt, search)
			args = append(args, filter.Term)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(likeStmt, " or")))
	}

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
