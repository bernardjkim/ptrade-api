package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/positions"
	"github.com/go-xorm/xorm"
)

// GetPositions will return the users current position for each stock
func GetPositions(DB *xorm.Engine, id int64) (positionsList positions.Positions, err error) {

	rows, err := DB.QueryString("CALL get_positions(?)", id)
	if err != nil {
		return
	}

	positionsList = positions.Positions{UserID: id}

	for _, row := range rows {
		p := positions.Position{}
		p.StockID, _ = strconv.ParseInt(row["stock_id"], 10, 64)
		p.Date = parseDate(string(row["date"]))
		p.Shares, _ = strconv.ParseInt(row["shares"], 10, 64)
		positionsList.Positions = append(positionsList.Positions, p)
	}
	return
}

// GetShares will return the users current position for a specifed stock
func GetShares(DB *xorm.Engine, id, stockID int64) (shares int64, err error) {

	rows, err := DB.Limit(1).QueryString("CALL get_shares(?, ?)", id, stockID)
	if err != nil {
		return
	}

	shares, _ = strconv.ParseInt(rows[0]["shares"], 10, 64)

	return
}
