package db

import (
	"context"
	"fmt"

	"github.com/Opsi/sparschwein/util"
	"github.com/jmoiron/sqlx"
)

func DoesTransactionExist(ctx context.Context, db sqlx.QueryerContext, transaction CreateTransaction) (bool, error) {
	// check if the transaction exists
	var selected Transaction
	query := `
		SELECT * FROM transactions
		WHERE from_holder_id = $1
		AND to_holder_id = $2
		AND amount = $3
		AND timestamp = $4`
	rows, err := db.QueryxContext(ctx, query,
		transaction.FromHolderID, transaction.ToHolderID, transaction.AmountInCents, transaction.Timestamp)
	if err != nil {
		return false, fmt.Errorf("select transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&selected)
		if err != nil {
			return false, fmt.Errorf("struct scan: %w", err)
		}

		isDataEqual, err := util.CompareNullJSONText(selected.Data, transaction.Data)
		if err != nil {
			return false, fmt.Errorf("compare null json text: %w", err)
		}
		if !isDataEqual {
			continue
		}
		return true, nil
	}
	return false, nil
}

func InsertTransaction(ctx context.Context, dbConn sqlx.ExtContext, create CreateTransaction) (*Transaction, error) {
	// check if the transaction exists
	exists, err := DoesTransactionExist(ctx, dbConn, create)
	if err != nil {
		return nil, fmt.Errorf("does transaction exist: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("transaction already exists")
	}

	// insert the transaction
	query := `
		INSERT INTO transactions
			(from_holder_id, to_holder_id, amount, timestamp, data)
	        VALUES (:from_holder_id, :to_holder_id, :amount, :timestamp, :data)
			RETURNING *`
	rows, err := sqlx.NamedQueryContext(ctx, dbConn, query, create)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}
	defer rows.Close()

	// Assuming your INSERT statement returns the newly created transaction
	if !rows.Next() {
		return nil, fmt.Errorf("no transaction returned")
	}
	var transaction Transaction
	err = rows.StructScan(&transaction)
	if err != nil {
		return nil, fmt.Errorf("struct scan: %w", err)
	}
	return &transaction, nil
}
