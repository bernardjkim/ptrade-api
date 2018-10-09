package db

import (
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/transfers"
	"github.com/go-xorm/xorm"
)

// GetTransfers will return all transfer orders made by the given user id
func GetTransfers(DB *xorm.Engine, id int64) (orders []transfers.TransferOrder) {

	rows, err := DB.QueryString("CALL get_transfer_orders(?)", id)
	checkError(err)

	for _, row := range rows {
		t := transfers.TransferOrder{}
		t.ID, _ = strconv.ParseInt(row["id"], 10, 64)
		t.UserID, _ = strconv.ParseInt(row["user_id"], 10, 64)
		t.DateStart = parseDate(string(row["date_start"]))
		t.DateEnd = parseDate(string(row["date_end"]))
		t.Balance, _ = strconv.ParseFloat(row["balance"], 64)
		orders = append(orders, t)
	}
	return
}

// NewTransfer will create a new transfer order for the specified user id
func NewTransfer(DB *xorm.Engine, id int64, balance float64) {
	_, err := DB.Exec("CALL new_transfer_order(?, ?)", id, balance)
	checkError(err)
	return
}

// GetTransfersTotal will return the total sum of balance transfers
func GetTransfersTotal(DB *xorm.Engine, id int64) (total float64) {

	rows, err := DB.QueryString("CALL get_transfers_total(?)", id)
	checkError(err)

	total, _ = strconv.ParseFloat(rows[0]["total"], 64)
	return
}
