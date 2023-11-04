package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Opsi/sparschwein/db"
	"github.com/jmoiron/sqlx"
)

type TransactionToCreate struct {
	Transaction    db.BaseTransaction
	FromIdentifier db.HolderIdentifier
	ToIdentifier   db.HolderIdentifier
}

type DryRunResult struct {
	ExistingHolders map[db.HolderIdentifier]db.Holder
	HoldersToCreate map[db.HolderIdentifier]db.CreateHolder
	Transactions    []TransactionToCreate
}

var _ json.Marshaler = DryRunResult{}

func (r DryRunResult) MarshalJSON() ([]byte, error) {
	asJson := struct {
		Holders      []db.CreateHolder
		Transactions []TransactionToCreate
	}{
		Holders:      make([]db.CreateHolder, 0, len(r.HoldersToCreate)),
		Transactions: r.Transactions,
	}
	for _, holder := range r.HoldersToCreate {
		asJson.Holders = append(asJson.Holders, holder)
	}
	return json.Marshal(asJson)
}

func (r DryRunResult) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("existingHolders", len(r.ExistingHolders)),
		slog.Int("holdersToCreate", len(r.HoldersToCreate)),
		slog.Int("transactions", len(r.Transactions)),
	)
}

func (r *DryRunResult) CheckHolder(ctx context.Context, dbConn sqlx.QueryerContext, cHolder db.CreateHolder) error {
	if _, ok := r.ExistingHolders[cHolder.HolderIdentifier]; ok {
		return nil
	}
	if _, ok := r.HoldersToCreate[cHolder.HolderIdentifier]; ok {
		return nil
	}

	// check if the holder exists
	holder, ok, err := db.GetHolderByIdentifier(ctx, dbConn, cHolder.HolderIdentifier)
	if err != nil {
		return fmt.Errorf("get holder: %w", err)
	}
	if ok {
		r.ExistingHolders[cHolder.HolderIdentifier] = *holder
		return nil
	}

	// if the holder doesn't exist, we need to create it
	r.HoldersToCreate[cHolder.HolderIdentifier] = cHolder
	return nil
}

func (r *DryRunResult) InsertHolders(ctx context.Context, dbConn sqlx.ExtContext) error {
	for cIdentifier, cHolder := range r.HoldersToCreate {
		newHolder, err := db.InsertHolder(ctx, dbConn, cHolder)
		if _, ok := r.ExistingHolders[cIdentifier]; ok {
			// this should never happen
			return fmt.Errorf("holder already exists")
		}
		if err != nil {
			return fmt.Errorf("insert holder: %w", err)
		}
		r.ExistingHolders[cIdentifier] = *newHolder
	}
	clear(r.HoldersToCreate)
	return nil
}

func DryRun(ctx context.Context,
	dbConn sqlx.QueryerContext,
	creators []TransactionCreator) (*DryRunResult, error) {
	// this is a dry run, so we just print the transactions
	// and holders that would be created

	result := &DryRunResult{
		ExistingHolders: make(map[db.HolderIdentifier]db.Holder),
		HoldersToCreate: make(map[db.HolderIdentifier]db.CreateHolder),
		Transactions:    make([]TransactionToCreate, 0),
	}

	// first we go over the holders and check which ones already exist
	for _, creator := range creators {
		if err := result.CheckHolder(ctx, dbConn, creator.FromHolder()); err != nil {
			return nil, fmt.Errorf("check from holder: %w", err)
		}
		if err := result.CheckHolder(ctx, dbConn, creator.ToHolder()); err != nil {
			return nil, fmt.Errorf("check to holder: %w", err)
		}
	}

	// then we go over the transactions and check which ones already exist
	for _, creator := range creators {
		createTransaction := TransactionToCreate{
			Transaction:    creator.Transaction(),
			FromIdentifier: creator.FromHolder().HolderIdentifier,
			ToIdentifier:   creator.ToHolder().HolderIdentifier,
		}
		// if the from holder doesn't exist, the transaction can't exist
		fromHolder, ok := result.ExistingHolders[createTransaction.FromIdentifier]
		if !ok {
			result.Transactions = append(result.Transactions, createTransaction)
			continue
		}
		// if the to holder doesn't exist, the transaction can't exist
		toHolder, ok := result.ExistingHolders[createTransaction.ToIdentifier]
		if !ok {
			result.Transactions = append(result.Transactions, createTransaction)
			continue
		}
		// if both holders exist, we need to check if the transaction exists
		ok, err := db.DoesTransactionExist(ctx, dbConn, db.CreateTransaction{
			BaseTransaction: createTransaction.Transaction,
			FromHolderID:    fromHolder.ID,
			ToHolderID:      toHolder.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("does transaction exist: %w", err)
		}
		if !ok {
			result.Transactions = append(result.Transactions, createTransaction)
		}
	}
	return result, nil
}
