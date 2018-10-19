package balances

import (
	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/db"
)

// Balances is a list of Balance objects
type Balances []Balance

// Balance represents an account balance for a user
type Balance db.BalanceTable

// TableName simply returns the table name
func (b *Balance) TableName() string {
	return "balances"
}
