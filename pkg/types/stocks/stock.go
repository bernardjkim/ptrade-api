package stocks

import "github.com/bkim0128/bjstock-rest-service/pkg/db"

type Stock db.StockTable

type StockList []Stock

func (s *Stock) TableName() string {
	return "stock_table"
}
