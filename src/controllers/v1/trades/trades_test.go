package trades_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bernardjkim/ptrade-api/pkg/types/trades"
	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/trades"
	"github.com/gorilla/mux"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	tradeHandler TradeHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	tradeHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = tradeHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = tradeHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = tradeHandler.DB.Exec("ALTER TABLE orders AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = tradeHandler.DB.Exec("ALTER TABLE trade_orders AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestGetTradesEmptyTable will test get trades with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/trades", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, tradeHandler.GetTrades)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

// TestGetTrades
func TestGetTrades(t *testing.T) {
	testSetup()

	// Test trades initialization
	tradeHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")
	tradeHandler.DB.Exec("CALL new_trade_order(1, 1, 5)")

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/trades", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, tradeHandler.GetTrades)

	Test.Equals(t, http.StatusOK, rr.Code)

	trade := trades.TradeOrder{OrderID: 1, StockID: 1, Shares: 5, Status: "FULFILLED"}

	exp := trades.TradeOrders{UserID: 1, Trades: []trades.TradeOrder{trade}}
	act := trades.TradeOrders{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, exp.UserID, act.UserID)
	for index, _ := range exp.Trades {
		Test.Equals(t, exp.Trades[index].OrderID, act.Trades[index].OrderID)
		Test.Equals(t, exp.Trades[index].StockID, act.Trades[index].StockID)
		Test.Equals(t, exp.Trades[index].Shares, act.Trades[index].Shares)
		Test.Equals(t, exp.Trades[index].Status, act.Trades[index].Status)
	}
}

// TestCreateTrade
func TestCreateTrades(t *testing.T) {
	testSetup()

	// Test trades initialization
	tradeHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")
	tradeHandler.DB.Exec("INSERT INTO stocks (symbol, name, price_per_share) VALUES ('AAPL', 'APPLE', 200)")

	req := httptest.NewRequest("POST", "/v1/users/{ID:[0-9]+}/trades", nil)

	req.Form = url.Values{
		"symbol": {"AAPL"},
		"shares": {"5"},
	}

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)

	// set user id in context
	ctx := context.WithValue(req.Context(), Users.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := Test.HandleRequest(req, tradeHandler.CreateTrade)

	Test.Equals(t, http.StatusCreated, rr.Code)

	trade := trades.TradeOrder{OrderID: 1, StockID: 1, Shares: 5, Status: "FULFILLED"}

	exp := trades.TradeOrders{UserID: 1, Trades: []trades.TradeOrder{trade}}
	act := trades.TradeOrders{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, exp.UserID, act.UserID)
	for index, _ := range exp.Trades {
		Test.Equals(t, exp.Trades[index].OrderID, act.Trades[index].OrderID)
		Test.Equals(t, exp.Trades[index].StockID, act.Trades[index].StockID)
		Test.Equals(t, exp.Trades[index].Shares, act.Trades[index].Shares)
		Test.Equals(t, exp.Trades[index].Status, act.Trades[index].Status)
	}
}

// TestInvalidSymbol
func TestInvalidSymbol(t *testing.T) {
	testSetup()

	// Test trades initialization
	tradeHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")

	req := httptest.NewRequest("POST", "/v1/users/{ID:[0-9]+}/trades", nil)

	req.Form = url.Values{
		"symbol": {"A"},
		"shares": {"5"},
	}

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)

	// set user id in context
	ctx := context.WithValue(req.Context(), Users.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := Test.HandleRequest(req, tradeHandler.CreateTrade)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided symbol id does not exist in database"
	act := strings.TrimSuffix(rr.Body.String(), "\n")

	Test.Equals(t, exp, act)
}
