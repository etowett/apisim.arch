package models

import (
	"apisim/app/db"
	"apisim/app/helpers"
	"context"
	"fmt"
	"strings"

	null "gopkg.in/guregu/null.v4"
)

const (
	createRecipientSQL         = `insert into recipients (message_id, phone, api_id, route, cost, currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`
	selectRecipientSQL         = `select r.id, r.message_id, r.phone, d.status, d.reason, r.api_id, r.route, r.cost, r.currency, r.created_at from recipients r left join dlrs d on r.id = d.recipient_id`
	selectMessageRecipientsSQL = selectRecipientSQL + ` where message_id=$1`
	countRecipientSQL          = `select count(r.id) from recipients r`
)

type (
	Recipient struct {
		SequentialIdentifier
		MessageID  int64       `json:"message_id"`
		Phone      string      `json:"phone"`
		Status     null.String `json:"status"`
		Reason     null.String `json:"reason"`
		Route      string      `json:"route"`
		Cost       string      `json:"cost"`
		Currency   string      `json:"currency"`
		Correlator string      `json:"correlator"`
		Timestamps
	}
)

func (r *Recipient) CountForMessage(
	ctx context.Context,
	db db.SQLOperations,
	messageID int64,
	filter *Filter,
) (int, error) {
	query, args := r.buildQuery(
		countRecipientSQL,
		messageID,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (r *Recipient) ForMessage(
	ctx context.Context,
	db db.SQLOperations,
	messageID int64,
	filter *Filter,
) ([]*Recipient, error) {
	recipients := make([]*Recipient, 0)
	query, args := r.buildQuery(
		selectRecipientSQL,
		messageID,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return recipients, err
	}

	for rows.Next() {
		var recipient Recipient
		err = rows.Scan(
			&recipient.ID,
			&recipient.MessageID,
			&recipient.Phone,
			&recipient.Status,
			&recipient.Reason,
			&recipient.Correlator,
			&recipient.Route,
			&recipient.Cost,
			&recipient.Currency,
			&recipient.CreatedAt,
		)
		if err != nil {
			return recipients, err
		}
		recipients = append(recipients, &recipient)
	}

	return recipients, err
}

func (r *Recipient) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	r.Timestamps.Touch()
	err := db.QueryRowContext(
		ctx,
		createRecipientSQL,
		r.MessageID,
		r.Phone,
		r.Correlator,
		r.Route,
		r.Cost,
		r.Currency,
		r.Timestamps.CreatedAt,
	).Scan(&r.ID)
	return err
}

func (r *Recipient) buildQuery(
	query string,
	messageID int64,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	conditions = append(conditions, fmt.Sprintf(" r.message_id=$%d", placeholder.Touch()))
	args = append(args, messageID)

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"phone", "message"}
		for _, col := range columns {
			search := fmt.Sprintf(" (lower(%s) like '%%' || $%d || '%%')", col, placeholder.Touch())
			likeStmt = append(likeStmt, search)
			args = append(args, filter.Term)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(likeStmt, " or")))
	}

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}
