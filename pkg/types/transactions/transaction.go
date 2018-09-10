package transactions

import (
	"github.com/bkim0128/bjstock-rest-service/pkg/db"
)

// type Portfolio []Transaction

type Transaction db.TransactionTable

func (t *Transaction) TableName() string {
	return "transaction_table"
}
