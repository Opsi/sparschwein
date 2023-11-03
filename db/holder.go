package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetHolderByIdentifier(ctx context.Context, db sqlx.QueryerContext, identifier HolderIdentifier) (Holder, error) {
	var holder Holder
	err := sqlx.GetContext(ctx,
		db,
		&holder,
		`SELECT * FROM holders WHERE type = $1 AND identifier = $2`,
		identifier.Type,
		identifier.Identifier)
	if err != nil {
		return Holder{}, fmt.Errorf("get holder by identifier: %w", err)
	}
	return holder, nil
}

func InsertHolder(ctx context.Context, db sqlx.ExtContext, createHolder CreateHolder) (Holder, error) {
	query := `
		INSERT INTO holders
			(type, identifier, name, parent_holder_id, data, favorite)
	        VALUES (:type, :identifier, :name, :parent_holder_id, :data, :favorite)
			RETURNING *`
	rows, err := sqlx.NamedQueryContext(ctx, db, query, createHolder)
	if err != nil {
		return Holder{}, err
	}

	// Assuming your INSERT statement returns the newly created holder
	if !rows.Next() {
		return Holder{}, fmt.Errorf("no holder returned")
	}
	var holder Holder
	err = rows.StructScan(&holder)
	if err != nil {
		return Holder{}, fmt.Errorf("struct scan: %w", err)
	}
	return holder, nil
}
