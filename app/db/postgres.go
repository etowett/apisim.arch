package db

import (
	"context"
	"database/sql"

	"github.com/revel/revel"

	// postgresql driver
	_ "github.com/lib/pq"
)

var db DBOperations

type SQLOperations interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type DBOperations interface {
	SQLOperations
	Begin() (*sql.Tx, error)
	Close() error
	Ping() error
	InTransaction(context.Context, func(context.Context, SQLOperations) error) error
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type AppDB struct {
	*sql.DB
}

func DB() DBOperations {
	if db == nil {
		revel.AppLog.Errorf("Database has not been initialized")
	}

	return db
}

func (db *AppDB) InTransaction(
	ctx context.Context,
	operations func(context.Context, SQLOperations) error,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err = operations(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}

func InitDB() {
	InitDBWithURL(
		revel.Config.StringDefault("db.url", ""),
	)
	revel.AppLog.Infof("DB Connected Successfully")
}

func InitDBWithURL(databaseURL string) DBOperations {
	appDB := newPostgresDBWithURL(databaseURL)
	db = &AppDB{appDB}

	err := db.Ping()
	if err != nil {
		revel.AppLog.Errorf("db ping failed because err=[%v]", err)
	}

	return db
}

func newPostgresDBWithURL(databaseURL string) *sql.DB {
	if databaseURL == "" {
		revel.AppLog.Errorf("database url is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		revel.AppLog.Errorf("sql.Open failed because err=[%v]", err)
	}

	return db
}
