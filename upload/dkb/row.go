package dkb

import (
	"time"

	"github.com/Opsi/sparschwein/db"
	"github.com/Opsi/sparschwein/upload"
)

type csvRow struct {
	BookingDate       time.Time
	ValueDate         time.Time
	Status            string
	Payer             string
	Payee             string
	Purpose           string
	TransactionType   string
	AmountInCents     int
	CreditorID        string
	MandateReference  string
	CustomerReference string
}

type transactionCreator struct {
	Row  csvRow
	Info *baseInfo
}

var _ upload.TransactionCreator = transactionCreator{}

func (t transactionCreator) fromIdentifier() string {
	// TODO
	return t.Row.Payer
}

func (t transactionCreator) toIdentifier() string {
	// TODO
	return t.Row.Payee
}

func (t transactionCreator) amountInCents() uint {
	if t.Row.AmountInCents < 0 {
		return uint(-t.Row.AmountInCents)
	}
	return uint(t.Row.AmountInCents)
}

func (t transactionCreator) Transaction() db.CreateTransaction {
	return db.CreateTransaction{
		FromIdentifier: t.fromIdentifier(),
		ToIdentifier:   t.toIdentifier(),
		AmountInCents:  t.amountInCents(),
		Time:           t.Row.ValueDate,
		Data:           t.Row,
	}
}

func (t transactionCreator) FromHolder() db.CreateHolder {
	// TODO
	return db.CreateHolder{
		Identifier: t.fromIdentifier(),
		Name:       t.Row.Payer,
		Data:       t,
	}
}

func (t transactionCreator) ToHolder() db.CreateHolder {
	// TODO
	return db.CreateHolder{
		Identifier: t.toIdentifier(),
		Name:       t.Row.Payee,
		Data:       t,
	}
}
