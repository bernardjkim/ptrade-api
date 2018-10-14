package db

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-xorm/xorm"
)

// GetPricePerShare will return pps for the specified stock
// func GetPricePerShare(DB *xorm.Engine, symbol string) (pps float64, err error) {
// 	return
// }

// SetPricePerShare will set the current pps for the specified stock
func SetPricePerShare(DB *xorm.Engine, symbol string) (err error) {

	// get current price for a share
	resp, err := http.Get("https://api.iextrading.com/1.0/stock/" + symbol + "/price")
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// TODO: format price to two dicimal places?
	price, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return
	}

	stockID, err := DB.Limit(1).QueryString("CALL get_id_from_symbol(?)", symbol)
	if err != nil {
		return
	}

	_, err = DB.Exec("CALL update_pps(?, ?)", stockID[0]["id"], price)
	if err != nil {
		return
	}

	return
}
