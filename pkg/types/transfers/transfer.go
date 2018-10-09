package transfers

import (
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// Transfers is a list of Transfer objects
type Transfers []Transfer

// Transfer represents a transfer order made by a user
type Transfer db.TransferOrderTable

// TransferOrder represents a transfer joined with order made by a user
type TransferOrder struct {
	ID        int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID    int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	DateStart time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd   time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	Balance   float64   `xorm:"FLOAT NOT NULL 'balance'" json:"balance" schema:"balance"`
}

// TableName simply returns the table name
func (t *Transfer) TableName() string {
	return "transfer_orders"
}
