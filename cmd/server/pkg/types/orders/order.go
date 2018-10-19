package orders

import (
	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/db"
)

// Orders is a list of order objects
type Orders []Orders

// Order represents a order made by a user
type Order db.OrderTable

// TableName simply returns the table name
func (o *Order) TableName() string {
	return "orders"
}
