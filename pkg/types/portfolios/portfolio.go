package portfolios

import (
	"time"
)

// PortfolioHistory represents a users portfolio value over a period of time
type PortfolioHistory struct {
	UserID  int64
	History []PortfolioValue
}

// PortfolioValue represents a users portfolio value at a certain date
type PortfolioValue struct {
	ID    int64
	Date  time.Time
	Value float64
}

// TableName simply returns the table name
func (p *PortfolioValue) TableName() string {
	return "portfolio_history"
}
