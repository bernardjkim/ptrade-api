package db

import (
	"strconv"

	"github.com/go-xorm/xorm"
)

// GetBalance will return the users current account balance
func GetBalance(DB *xorm.Engine, id int64) (balance float64) {
	rows, err := DB.Limit(1).QueryString("CALL get_balance(?)", id)
	checkError(err)

	balance, _ = strconv.ParseFloat(rows[0]["balance"], 64)
	return
}
