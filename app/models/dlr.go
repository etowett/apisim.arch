package models

import (
	"apisim/app/db"
	"context"
	"time"
)

const (
	createDlrSQL = `insert into dlrs (recipient_id, status, reason, received_at, created_at) VALUES ($1, $2, $3, $4, $5) returning id`
)

type (
	Dlr struct {
		SequentialIdentifier
		RecipientID int64     `json:"recipient_id"`
		Status      string    `json:"status"`
		Reason      string    `json:"reason"`
		ReceivedAt  time.Time `json:"received_at"`
		Timestamps
	}
)

func (d *Dlr) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	d.Timestamps.Touch()
	err := db.QueryRowContext(
		ctx,
		createDlrSQL,
		d.RecipientID,
		d.Status,
		d.Reason,
		d.ReceivedAt,
		d.Timestamps.CreatedAt,
	).Scan(&d.ID)
	return err
}
