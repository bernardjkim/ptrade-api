package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/trades"
	"github.com/go-xorm/xorm"
)

// GetTrades will return all transfer orders made by the given user id
func GetTrades(DB *xorm.Engine, id int64) (orders trades.TradeOrders, err error) {

	rows, err := DB.QueryString("CALL get_trade_orders(?)", id)
	if err != nil {
		return
	}

	orders = trades.TradeOrders{UserID: id}

	for _, row := range rows {
		t := trades.TradeOrder{}
		t.OrderID, _ = strconv.ParseInt(row["id"], 10, 64)
		t.DateStart = parseDate(string(row["date_start"]))
		t.DateEnd = parseDate(string(row["date_end"]))
		t.StockID, _ = strconv.ParseInt(row["stock_id"], 10, 64)
		t.Shares, _ = strconv.ParseInt(row["shares"], 10, 64)
		t.PricePerShare, _ = strconv.ParseFloat(row["price_per_share"], 64)
		t.Status, _ = row["status"]
		orders.Trades = append(orders.Trades, t)
	}
	return
}

// NewTrade will create a new trade order for the specified user id
func NewTrade(DB *xorm.Engine, userID int64, stockID int64, shares int64) (err error) {
	_, err = DB.Exec("CALL new_trade_order(?, ?, ?)", userID, stockID, shares)
	return
}
