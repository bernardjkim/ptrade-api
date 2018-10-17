package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/portfolios"
	"github.com/go-xorm/xorm"
)

// GetPortfolioHistory will return history of user portfolio value
func GetPortfolioHistory(DB *xorm.Engine, id int64) (history portfolios.PortfolioHistory, err error) {

	rows, err := DB.QueryString("CALL get_portfolio_history(?)", id)
	if err != nil {
		return
	}

	history.UserID = id
	history.History = []portfolios.PortfolioValue{}

	for _, row := range rows {
		pv := portfolios.PortfolioValue{}
		pv.ID, _ = strconv.ParseInt(row["id"], 10, 64)
		pv.Date = parseDate(string(row["date"]))
		pv.Value, _ = strconv.ParseFloat(row["value"], 64)

		history.History = append(history.History, pv)
	}
	return
}

// GetProfit will return profit
func GetProfit(DB *xorm.Engine, id int64) (profit float64) {

	rows, err := DB.QueryString("CALL get_profit(?)", id)
	if err != nil {
		return
	}

	profit, _ = strconv.ParseFloat(rows[0]["profit"], 64)

	return
}
