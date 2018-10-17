package portfolios

import (
	"time"
)

// PortfolioHistory represents a users portfolio value over a period of time
type PortfolioHistory struct {
	UserID  int64            `json:"user_id"`
	History []PortfolioValue `json:"history"`
}

// PortfolioValue represents a users portfolio value at a certain date
type PortfolioValue struct {
	ID    int64     `json:"id"`
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// TableName simply returns the table name
func (p *PortfolioValue) TableName() string {
	return "portfolio_history"
}
