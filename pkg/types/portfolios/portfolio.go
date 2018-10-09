package portfolios

import (
	"github.com/bernardjkim/ptrade-api/pkg/db"
)

// PortfolioHistory represents a users portfolio value over a period of time
type PortfolioHistory []PortfolioValue

// PortfolioValue represents a users portfolio value at a certain date
type PortfolioValue db.PortfolioHistoryTable

// TableName simply returns the table name
func (p *PortfolioValue) TableName() string {
	return "portfolio_history"
}
