package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/balances"
	"github.com/go-xorm/xorm"
)

// GetBalance will return the users current account balance
func GetBalance(DB *xorm.Engine, id int64) (balance balances.Balance, err error) {
	rows, err := DB.Limit(1).QueryString("CALL get_balance(?)", id)
	if err != nil || len(rows) < 1 {
		return
	}

	balance.ID, _ = strconv.ParseInt(rows[0]["id"], 10, 64)
	balance.UserID, _ = strconv.ParseInt(rows[0]["user_id"], 10, 64)
	balance.DateStart = parseDate(string(rows[0]["date_start"]))
	balance.DateEnd = parseDate(string(rows[0]["date_end"]))
	balance.Balance, _ = strconv.ParseFloat(rows[0]["balance"], 64)
	return
}
