package db

import "time"

type CreateLedger struct {
	Identifier string
	Name       string
	Data       any
}

type Ledger struct{}

type CreateTransaction struct {
	FromIdentifier string
	ToIdentifier   string
	AmountInCents  uint
	Time           time.Time
	Data           any
}

type Transaction struct{}

type Tag struct{}
