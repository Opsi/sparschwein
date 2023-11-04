package upload

import (
	"github.com/Opsi/sparschwein/db"
)

type TransactionCreator interface {
	Transaction() db.BaseTransaction
	FromHolder() db.CreateHolder
	ToHolder() db.CreateHolder
}
