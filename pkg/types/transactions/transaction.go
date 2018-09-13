package transactions

import (
	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// type Portfolio []Transaction

type Transaction db.TransactionTable

func (t *Transaction) TableName() string {
	return "transaction"
}
