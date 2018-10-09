package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/positions"
	"github.com/go-xorm/xorm"
)

// GetPositions will return the users current position for each stock
func GetPositions(DB *xorm.Engine, id int64) (positionsList positions.Positions) {

	rows, err := DB.QueryString("CALL get_positions(?)", id)
	checkError(err)

	for _, row := range rows {
		p := positions.Position{}
		p.ID, _ = strconv.ParseInt(row["id"], 10, 64)
		p.UserID, _ = strconv.ParseInt(row["user_id"], 10, 64)
		p.DateStart = parseDate(string(row["date_start"]))
		p.DateStart = parseDate(string(row["date_end"]))
		p.StockID, _ = strconv.ParseInt(row["stock_id"], 10, 64)
		p.Shares, _ = strconv.ParseInt(row["shares"], 10, 64)
		positionsList = append(positionsList, p)
	}
	return
}

// GetShares will return the users current position for a specifed stock
func GetShares(DB *xorm.Engine, id, stockID int64) (shares int64) {

	rows, err := DB.Limit(1).QueryString("CALL get_shares(?, ?)", id, stockID)
	checkError(err)

	shares, _ = strconv.ParseInt(rows[0]["shares"], 10, 64)

	return
}
