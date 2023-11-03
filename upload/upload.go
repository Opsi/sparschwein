package upload

import (
	"github.com/Opsi/sparschwein/db"
)

type TransactionCreator interface {
	Transaction() db.CreateTransaction
	FromLedger() db.CreateLedger
	ToLedger() db.CreateLedger
}
