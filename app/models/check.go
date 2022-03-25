package models

import (
	"apisim/app/db"
	"context"
	"database/sql"
)

const (
	getDatabaseVersionSQL = `select version()`
	getDatabaseOneSQL     = `select 1`
)

type Checker struct{}

func (m *Checker) Version(
	ctx context.Context,
	db db.SQLOperations,
) (string, error) {
	row := db.QueryRowContext(ctx, getDatabaseVersionSQL)
	return m.scanRow(row)
}

func (m *Checker) One(
	ctx context.Context,
	db db.SQLOperations,
) (string, error) {
	row := db.QueryRowContext(ctx, getDatabaseOneSQL)
	return m.scanRow(row)
}

func (*Checker) scanRow(
	row *sql.Row,
) (string, error) {
	var val string
	err := row.Scan(&val)
	return val, err
}
