package trades

import (
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// Trades is a list of Trade objects
type Trades []Trade

// Trade represents a trade order made by a user
type Trade db.TradeOrderTable

// TradeOrders represent a list of transfer order made by a user
type TradeOrders struct {
	UserID int64
	Trades []TradeOrder
}

// TradeOrder represents a transfer joined with order made by a user
type TradeOrder struct {
	OrderID       int64     `xorm:"INTEGER NOT NULL 'order_id'" json:"order_id" schema:"order_id"`
	DateStart     time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd       time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	StockID       int64     `xorm:"INTEGER NOT NULL 'stock_id'" json:"stock_id" schema:"stock_id"`
	Shares        int64     `xorm:"INTEGER NOT NULL 'shares'" json:"shares" schema:"shares"`
	PricePerShare float64   `xorm:"FLOAT NOT NULL 'price_per_share'" json:"price_per_share" schema:"price_per_share"`
	Status        string    `xorm:"VARCHAR(20) 'status'" json:"status" schema:"status"`
}

// TableName simply returns the table name
func (t *Trade) TableName() string {
	return "trade_orders"
}
