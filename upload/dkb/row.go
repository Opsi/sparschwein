package dkb

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Opsi/sparschwein/db"
	"github.com/Opsi/sparschwein/upload"
	"github.com/jmoiron/sqlx/types"
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

var _ slog.LogValuer = csvRow{}

func (r csvRow) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Time("bookingDate", r.BookingDate),
		slog.Time("valueDate", r.ValueDate),
		slog.String("status", r.Status),
		slog.String("payer", r.Payer),
		slog.String("payee", r.Payee),
		slog.String("purpose", r.Purpose),
		slog.String("transactionType", r.TransactionType),
		slog.Int("amountInCents", r.AmountInCents),
		slog.String("creditorID", r.CreditorID),
		slog.String("mandateReference", r.MandateReference),
		slog.String("customerReference", r.CustomerReference),
	)
}

type transactionCreator struct {
	Row     csvRow
	Account *account
}

var _ upload.TransactionCreator = transactionCreator{}

func (t transactionCreator) Transaction() db.BaseTransaction {
	data, err := json.Marshal(t.Row)
	if err != nil {
		slog.Error("error parsing transaction data",
			slog.String("error", err.Error()),
			slog.Any("row", t.Row))
	}
	return db.BaseTransaction{
		AmountInCents: max(t.Row.AmountInCents, -t.Row.AmountInCents),
		Timestamp:     t.Row.ValueDate,
		Data: types.NullJSONText{
			JSONText: data,
			Valid:    true,
		},
		ParentTransactionID: nil,
	}
}

func (t transactionCreator) FromHolder() db.CreateHolder {
	if t.Row.AmountInCents < 0 {
		// The owner of the account is the payer
		return t.Account.createHolder()
	}
	// The owner of the account is the payee
	return db.CreateHolder{
		HolderIdentifier: db.HolderIdentifier{
			Type:       "dkb/payer",
			Identifier: t.Row.Payer,
		},
		ParentHolderID: nil,
		Favorite:       false,
		Name:           t.Row.Payer,
		Data: types.NullJSONText{
			JSONText: nil,
			Valid:    false,
		},
	}
}

func (t transactionCreator) ToHolder() db.CreateHolder {
	if t.Row.AmountInCents < 0 {
		// The owner of the account is the payer
		return db.CreateHolder{
			HolderIdentifier: db.HolderIdentifier{
				Type:       "dkb/payee",
				Identifier: t.Row.Payee,
			},
			ParentHolderID: nil,
			Favorite:       false,
			Name:           t.Row.Payee,
			Data: types.NullJSONText{
				JSONText: nil,
				Valid:    false,
			},
		}
	}
	// The owner of the account is the payee
	return t.Account.createHolder()
}
