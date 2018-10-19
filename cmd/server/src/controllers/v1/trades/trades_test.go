package trades_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/trades"
	Users "github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/users"
	Test "github.com/bernardjkim/ptrade-api/cmd/server/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/cmd/server/src/controllers/v1/trades"
	"github.com/gorilla/mux"
)

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
	if _, err = tradeHandler.DB.Exec("DELETE FROM users"); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = tradeHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1"); err != nil {
		log.Fatal(err.Error())
	}

	if _, err = tradeHandler.DB.Exec("DELETE FROM stocks"); err != nil {
		log.Fatal(err.Error())
	}
	if _, err = tradeHandler.DB.Exec("ALTER TABLE stocks AUTO_INCREMENT = 1"); err != nil {
		log.Fatal(err.Error())
	}

	if _, err = tradeHandler.DB.Exec("ALTER TABLE orders AUTO_INCREMENT = 1"); err != nil {
		log.Fatal(err.Error())
	}

	if _, err = tradeHandler.DB.Exec("ALTER TABLE trade_orders AUTO_INCREMENT = 1"); err != nil {
		log.Fatal(err.Error())
	}
}

// getTradesReq returns a new trade req with the given fields
func getTradesReq(id string) (req *http.Request) {
	req = httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/trades", nil)
	// set mux vars
	vars := map[string]string{
		"ID": id,
	}
	req = mux.SetURLVars(req, vars)
	return
}

// createTradeReq returns a new trade req with the given fields
func createTradeReq(id, symbol, shares string) (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/users/{ID:[0-9]+}/trades", nil)

	req.Form = url.Values{
		"symbol": {symbol},
		"shares": {shares},
	}

	// set mux vars
	vars := map[string]string{
		"ID": id,
	}
	req = mux.SetURLVars(req, vars)

	// set user id in context
	ctx := context.WithValue(req.Context(), Users.UserIDKey, 1)
	req = req.WithContext(ctx)
	return
}

// tradeOrder returns a new TradeOrder with the given fields
func tradeOrder(orderID, stockID, shares int64, status string) (trade trades.TradeOrder) {
	trade = trades.TradeOrder{OrderID: orderID, StockID: stockID, Shares: shares, Status: status}
	return
}

// compareTrades will compare the fields of the exp and act
func compareTrades(t *testing.T, exp trades.TradeOrders, act trades.TradeOrders) {
	Test.Equals(t, fmt.Sprintf("%s: %d", "user id", exp.UserID), fmt.Sprintf("%s: %d", "user id", act.UserID))
	for index, _ := range exp.Trades {
		Test.Equals(t, fmt.Sprintf("%s: %d", "order id", exp.Trades[index].OrderID), fmt.Sprintf("%s: %d", "order id", act.Trades[index].OrderID))
		Test.Equals(t, fmt.Sprintf("%s: %d", "stock id", exp.Trades[index].StockID), fmt.Sprintf("%s: %d", "stock id", act.Trades[index].StockID))
		Test.Equals(t, fmt.Sprintf("%s: %d", "shares", exp.Trades[index].Shares), fmt.Sprintf("%s: %d", "shares", act.Trades[index].Shares))
		Test.Equals(t, fmt.Sprintf("%s: %s", "status", exp.Trades[index].Status), fmt.Sprintf("%s: %s", "status", act.Trades[index].Status))
	}
}

// TestGetTradesEmptyTable will test get trades with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := getTradesReq("1")
	rr := Test.HandleRequest(req, tradeHandler.GetTrades)
	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}

// TestGetTrades
func TestGetTrades(t *testing.T) {
	testSetup()
	Test.NewUser("test1", "test1", "test1", "test1")
	Test.NewStock("AAPL", "APPLE", 200)
	req := createTradeReq("1", "AAPL", "5")
	_ = Test.HandleRequest(req, tradeHandler.CreateTrade)

	req = getTradesReq("1")
	rr := Test.HandleRequest(req, tradeHandler.GetTrades)
	Test.Equals(t, http.StatusOK, rr.Code)

	tradeList := []trades.TradeOrder{}
	tradeList = append(tradeList, tradeOrder(1, 1, 5, "FULFILLED"))

	exp := trades.TradeOrders{UserID: 1, Trades: tradeList}
	act := trades.TradeOrders{}
	json.NewDecoder(rr.Body).Decode(&act)

	compareTrades(t, exp, act)
}

// TestCreateTrade will test creating trade orders
func TestCreateTrades(t *testing.T) {
	testSetup()
	Test.NewUser("test1", "test1", "test1", "test1")
	Test.NewStock("AAPL", "APPLE", 200)
	Test.NewStock("MSFT", "MICROSOFT", 100)
	Test.NewStock("NVDA", "NVIDIA", 250)

	// Test single entry
	req := createTradeReq("1", "AAPL", "5")
	rr := Test.HandleRequest(req, tradeHandler.CreateTrade)
	Test.Equals(t, http.StatusCreated, rr.Code)

	tradeList := []trades.TradeOrder{}
	tradeList = append(tradeList, tradeOrder(1, 1, 5, "FULFILLED"))

	exp := trades.TradeOrders{UserID: 1, Trades: tradeList}
	act := trades.TradeOrders{}
	json.NewDecoder(rr.Body).Decode(&act)
	compareTrades(t, exp, act)

	// Test multiple entries
	req = createTradeReq("1", "AAPL", "5")
	_ = Test.HandleRequest(req, tradeHandler.CreateTrade)
	req = createTradeReq("1", "MSFT", "5")
	_ = Test.HandleRequest(req, tradeHandler.CreateTrade)
	req = createTradeReq("1", "NVDA", "5")
	rr = Test.HandleRequest(req, tradeHandler.CreateTrade)
	Test.Equals(t, http.StatusCreated, rr.Code)

	tradeList = append(tradeList, tradeOrder(2, 1, 5, "FULFILLED"))
	tradeList = append(tradeList, tradeOrder(3, 2, 5, "FULFILLED"))
	tradeList = append(tradeList, tradeOrder(4, 3, 5, "FULFILLED"))

	exp = trades.TradeOrders{UserID: 1, Trades: tradeList}
	act = trades.TradeOrders{}
	json.NewDecoder(rr.Body).Decode(&act)
	compareTrades(t, exp, act)
}

// TestInvalidSymbol will test result of creating a trade order for an invalid symbol
func TestInvalidSymbol(t *testing.T) {
	testSetup()
	Test.NewUser("test1", "test1", "test1", "test1")

	req := createTradeReq("1", "AAPL", "5")
	rr := Test.HandleRequest(req, tradeHandler.CreateTrade)
	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided symbol id does not exist in database"
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}
