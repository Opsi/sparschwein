package db

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type HolderIdentifier struct {
	Type       string
	Identifier string
}

func (h HolderIdentifier) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", h.Type),
		slog.String("identifier", h.Identifier),
	)
}

type CreateHolder struct {
	HolderIdentifier
	Name           string
	ParentHolderID *int `db:"parent_holder_id"`
	Data           types.NullJSONText
	Favorite       bool
}

type Holder struct {
	CreateHolder
	ID        int
	CreatedAt time.Time `db:"created_at"`
}

type BaseTransaction struct {
	AmountInCents       int `db:"amount"`
	Timestamp           time.Time
	Data                types.NullJSONText
	ParentTransactionID *int `db:"parent_transaction_id"`
}

type CreateTransaction struct {
	BaseTransaction
	FromHolderID int `db:"from_holder_id"`
	ToHolderID   int `db:"to_holder_id"`
}

type Transaction struct {
	CreateTransaction
	ID        int
	CreatedAt time.Time `db:"created_at"`
}

type Tag struct {
	ID          int
	Name        string
	ParentTagID *int      `db:"parent_tag_id"`
	CreatedAt   time.Time `db:"created_at"`
}
