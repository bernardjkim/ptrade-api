package transfers

import (
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// Transfers is a list of Transfer objects
type Transfers []Transfer

// Transfer represents a transfer order made by a user
type Transfer db.TransferOrderTable

// TransferOrders represent a list of transfer order made by a user
type TransferOrders struct {
	UserID    int64           `json:"user_id"`
	Transfers []TransferOrder `json:"transfers"`
}

// TransferOrder represents a transfer joined with order made by a user
type TransferOrder struct {
	OrderID   int64     `xorm:"INTEGER NOT NULL 'order_id'" json:"order_id" schema:"order_id"`
	DateStart time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd   time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	Balance   float64   `xorm:"FLOAT NOT NULL 'balance'" json:"balance" schema:"balance"`
	Status    string    `xorm:"VARCHAR(20) 'status'" json:"status" schema:"status"`
}

// TableName simply returns the table name
func (t *Transfer) TableName() string {
	return "transfer_orders"
}
