package positions

import (
	"time"
)

// Positions is a list of Position objects
type Positions struct {
	UserID    int64
	Positions []Position
}

// Position represents a users position for a specific stock
type Position struct {
	StockID int64
	Date    time.Time
	Shares  int64
}

// TableName simply returns the table name
func (h *Position) TableName() string {
	return "positions"
}
