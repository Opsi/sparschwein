package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetHolderByIdentifier(ctx context.Context, db sqlx.QueryerContext, identifier HolderIdentifier) (*Holder, bool, error) {
	var holder Holder
	const query = "SELECT * FROM holders WHERE type = $1 AND identifier = $2"
	err := sqlx.GetContext(ctx, db, &holder, query, identifier.Type, identifier.Identifier)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("select holder: %w", err)
	}
	return &holder, true, nil
}

func InsertHolder(ctx context.Context, db sqlx.ExtContext, createHolder CreateHolder) (*Holder, error) {
	query := `
		INSERT INTO holders
			(type, identifier, name, parent_holder_id, data, favorite)
	        VALUES (:type, :identifier, :name, :parent_holder_id, :data, :favorite)
			RETURNING *`
	rows, err := sqlx.NamedQueryContext(ctx, db, query, createHolder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Assuming your INSERT statement returns the newly created holder
	if !rows.Next() {
		return nil, fmt.Errorf("no holder returned")
	}
	var holder Holder
	err = rows.StructScan(&holder)
	if err != nil {
		return nil, fmt.Errorf("struct scan: %w", err)
	}
	return &holder, nil
}

func GetHolders(ctx context.Context, db sqlx.QueryerContext) ([]Holder, error) {
	var holders []Holder
	const query = "SELECT * FROM holders ORDER BY favorite DESC, id ASC"
	err := sqlx.SelectContext(ctx, db, &holders, query)
	if err != nil {
		return nil, fmt.Errorf("select holders: %w", err)
	}
	return holders, nil
}
