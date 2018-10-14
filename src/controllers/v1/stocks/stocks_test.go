package stocks_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/stocks"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	stockHandler StockHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	stockHandler.Init(db)
}

func TestGetStocks(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/stocks", nil)
	rr := Test.HandleRequest(req, stockHandler.GetStocks)

	Test.Equals(t, http.StatusOK, rr.Code)

	// TODO: is it necessary to test response body?
}
