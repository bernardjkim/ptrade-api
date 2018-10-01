package stocktransactions

import (
	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// StockTransaction represents a stock transaction
type Transaction db.StockTransactionTable

// TableName returns table name
func (t *Transaction) TableName() string {
	return "stock_transactions"
}
