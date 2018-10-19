package stocks

import (
	"encoding/json"
	"log"
	"net/http"

	Stocks "github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/stocks"
	ORM "github.com/bernardjkim/ptrade-api/cmd/server/src/system/db"

	"github.com/go-xorm/xorm"
)

// StockHandler struct needs to be initialized with a database connection.
type StockHandler struct {
	DB *xorm.Engine
}

// Init function initializes db connection
func (s *StockHandler) Init(DB *xorm.Engine) {
	s.DB = DB
}

// GetStocks function will return a list of available stock data containing
// id, symbol, and company name.
func (s *StockHandler) GetStocks(w http.ResponseWriter, r *http.Request) {

	stockList := Stocks.StockList{}

	// get list of available stocks from database
	if err := ORM.Find(s.DB, &Stocks.Stock{}, &stockList); err != nil {
		log.Println(err)
		http.Error(w, "Unable to get stock list", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(stockList)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
