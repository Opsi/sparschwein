package dkb

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Opsi/sparschwein/db"
	"github.com/jmoiron/sqlx/types"
)

type account struct {
	HolderType string
	IBAN       string
}

func (a account) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("holderType", a.HolderType),
		slog.String("iban", a.IBAN),
	)
}

func (a account) holderIndentifier() db.HolderIdentifier {
	return db.HolderIdentifier{
		Type:       "iban",
		Identifier: a.IBAN,
	}
}

func (a account) createHolder() db.CreateHolder {
	// The owner of the account is the payer
	accountInfoBytes, err := json.Marshal(a)
	if err != nil {
		slog.Error("error parsing account info",
			slog.String("error", err.Error()),
			slog.Any("account", a))
	}
	return db.CreateHolder{
		HolderIdentifier: a.holderIndentifier(),
		ParentHolderID:   nil,
		Favorite:         true,
		Name:             fmt.Sprintf("%s %s", a.HolderType, a.IBAN),
		Data: types.NullJSONText{
			JSONText: accountInfoBytes,
			Valid:    true,
		},
	}
}
