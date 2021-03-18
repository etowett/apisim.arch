package models

import (
	"apisim/app/db"
	"apisim/app/helpers"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	createApiKeySQL             = `insert into api_keys (user_id, provider, name, access_id, access_secret_hash, dlr_url, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`
	selectApiKeySQL             = `select id, user_id, provider, name, access_id, access_secret_hash, dlr_url, created_at from api_keys`
	selectUserApiKeySQL         = selectApiKeySQL + ` where user_id=$1`
	getApiKeyByIDSQL            = selectApiKeySQL + ` where id=$1`
	getApiKeyByAccountIDSQL     = selectApiKeySQL + ` where access_id=$1`
	getApiKeyUserAndAccessIDSQL = selectApiKeySQL + " where access_id=$1 and user_id=$2"
	countApiKeySQL              = `select count(id) from api_keys`
	deleteApiKeySQL             = `delete from api_keys where id=$1`
	updateApiKeySQL             = `update api_keys set (name, dlr_url, updated_at) = ($1, $2, $3) where id=$4`
)

type (
	ApiKey struct {
		SequentialIdentifier
		UserID           int64  `json:"user_id"`
		Provider         string `json:"provider"`
		Name             string `json:"name"`
		AccessID         string `json:"access_id"`
		AccessSecretHash string `json:"-"`
		DlrURL           string `json:"dlr_url"`
		Timestamps
	}
)

func (a *ApiKey) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*ApiKey, error) {
	apiKeys := make([]*ApiKey, 0)

	query, args := a.buildQuery(
		selectApiKeySQL,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return apiKeys, err
	}

	for rows.Next() {
		var apiKey ApiKey
		err = rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.Provider,
			&apiKey.Name,
			&apiKey.AccessID,
			&apiKey.AccessSecretHash,
			&apiKey.DlrURL,
			&apiKey.CreatedAt,
		)
		if err != nil {
			return apiKeys, err
		}
		apiKeys = append(apiKeys, &apiKey)
	}

	return apiKeys, err
}

func (a *ApiKey) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*ApiKey, error) {
	var apiKey ApiKey
	row := db.QueryRowContext(ctx, getApiKeyByIDSQL, id)

	err := a.scan(row, &apiKey)
	return &apiKey, err
}

func (a *ApiKey) ByUserAndAccessID(
	ctx context.Context,
	db db.SQLOperations,
	accessID string,
	userID int64,
) (*ApiKey, error) {
	var apiKey ApiKey
	row := db.QueryRowContext(ctx, getApiKeyUserAndAccessIDSQL, accessID, userID)

	err := a.scan(row, &apiKey)
	return &apiKey, err
}

func (a *ApiKey) ByAccountID(
	ctx context.Context,
	db db.SQLOperations,
	accountID string,
) (*ApiKey, error) {
	var apiKey ApiKey
	row := db.QueryRowContext(ctx, getApiKeyByAccountIDSQL, accountID)

	err := a.scan(row, &apiKey)
	return &apiKey, err
}

func (*ApiKey) Delete(
	ctx context.Context,
	db db.SQLOperations,
	apiID int64,
) error {
	res, err := db.ExecContext(
		ctx,
		deleteApiKeySQL,
		apiID,
	)
	if err != nil {
		return fmt.Errorf("exec delete failed: %v", err)
	}
	_, err = res.RowsAffected()
	return err
}

func (a *ApiKey) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := a.buildQuery(
		countApiKeySQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (a *ApiKey) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	a.Timestamps.Touch()
	var err error
	if a.IsNew() {
		err = db.QueryRowContext(
			ctx,
			createApiKeySQL,
			a.UserID,
			a.Provider,
			a.Name,
			a.AccessID,
			a.AccessSecretHash,
			a.DlrURL,
			a.Timestamps.CreatedAt,
		).Scan(&a.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updateApiKeySQL,
		a.Name,
		a.DlrURL,
		a.Timestamps.UpdatedAt,
		a.ID,
	)
	return err
}

func (*ApiKey) scan(
	row *sql.Row,
	apiKey *ApiKey,
) error {
	return row.Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.Provider,
		&apiKey.Name,
		&apiKey.AccessID,
		&apiKey.AccessSecretHash,
		&apiKey.DlrURL,
		&apiKey.CreatedAt,
	)
}

func (a *ApiKey) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"name", "access_id"}
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
