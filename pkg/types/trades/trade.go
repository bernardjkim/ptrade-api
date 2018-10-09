package trades

import (
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// Trades is a list of Trade objects
type Trades []Trade

// Trade represents a trade order made by a user
type Trade db.TradeOrderTable

// TradeOrder represents a trade joined with order made by a user
type TradeOrder struct {
	ID            int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID        int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	DateStart     time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd       time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	StockID       int64     `xorm:"INTEGER NOT NULL 'stock_id'" json:"stock_id" schema:"stock_id"`
	Shares        int64     `xorm:"INTEGER NOT NULL 'shares'" json:"shares" schema:"shares"`
	PricePerShare float64   `xorm:"FLOAT NOT NULL 'price_per_share'" json:"price_per_share" schema:"price_per_share"`
}

// TableName simply returns the table name
func (t *Trade) TableName() string {
	return "trade_orders"
}
