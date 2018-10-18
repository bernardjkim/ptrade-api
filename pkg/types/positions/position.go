package positions

import (
	"time"
)

// Positions is a list of Position objects
type Positions struct {
	UserID    int64      `json:"user_id"`
	Positions []Position `json:"positions"`
}

// Position represents a users position for a specific stock
type Position struct {
	StockID       int64     `json:"stock_id"`
	Date          time.Time `json:"date"`
	Symbol        string    `json:"symbol"`
	PricePerShare float64   `json:"price_per_share"`
	Shares        int64     `json:"shares"`
}

// TableName simply returns the table name
func (h *Position) TableName() string {
	return "positions"
}
