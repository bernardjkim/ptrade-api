package stocks

import "github.com/bernardjkim/ptrade-api/pkg/db"

type Stock db.StockTable

type StockList []Stock

func (s *Stock) TableName() string {
	return "stocks"
}
