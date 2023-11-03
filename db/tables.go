package db

import "time"

type CreateHolder struct {
	Identifier string
	Name       string
	Data       any
}

type Holder struct{}

type CreateTransaction struct {
	FromIdentifier string
	ToIdentifier   string
	AmountInCents  uint
	Time           time.Time
	Data           any
}

type Transaction struct{}

type Tag struct{}
