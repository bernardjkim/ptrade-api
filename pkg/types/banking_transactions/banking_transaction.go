package bankingtransactions

import (
	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// StockTransaction represents a stock transaction
type Transaction db.BankingTransactionTable

// TableName returns table name
func (t *Transaction) TableName() string {
	return "banking_transactions"
}
