package db

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type HolderIdentifier struct {
	Type       string
	Identifier string
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
	AmountInCents       int `db:"amount_in_cents"`
	Timestamp           time.Time
	Data                types.NullJSONText
	ParentTransactionID *int `db:"parent_transaction_id"`
}

type CreateTransaction struct {
	BaseTransaction
	FromIdentifier HolderIdentifier
	ToIdentifier   HolderIdentifier
}

type Transaction struct {
	BaseTransaction
	ID           int
	FromHolderID int       `db:"from_holder_id"`
	ToHolderID   int       `db:"to_holder_id"`
	CreatedAt    time.Time `db:"created_at"`
}

type Tag struct {
	ID          int
	Name        string
	ParentTagID *int      `db:"parent_tag_id"`
	CreatedAt   time.Time `db:"created_at"`
}
