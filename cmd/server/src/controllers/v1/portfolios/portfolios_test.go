package portfolios_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/portfolios"
	. "github.com/bernardjkim/ptrade-api/cmd/server/src/controllers/v1/portfolios"
	Test "github.com/bernardjkim/ptrade-api/cmd/server/src/controllers/v1/test"
	"github.com/gorilla/mux"
)

var (
	portfolioHandler PortfolioHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	portfolioHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = portfolioHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = portfolioHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = portfolioHandler.DB.Exec("ALTER TABLE portfolio_history AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestEmptyTable will test get portfolio history endpoint on an empty table.
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/charts", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, portfolioHandler.GetPortfolioHistory)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}

// TestGetHistory will test getting portfolio history endpoint.
func TestGetHistory(t *testing.T) {
	testSetup()

	// Test Single Portfolio Value
	portfolioHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test')")
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (1, '2018-10-10 21:05:29', 123)")

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/charts", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, portfolioHandler.GetPortfolioHistory)

	Test.Equals(t, http.StatusOK, rr.Code)

	date, _ := time.Parse("2006-01-02T15:04:05Z", "2018-10-10T21:05:29Z")
	exp := portfolios.PortfolioHistory{
		UserID: 1,
		History: []portfolios.PortfolioValue{
			{ID: 1, Date: date, Value: 123},
		},
	}
	act := portfolios.PortfolioHistory{}
	json.NewDecoder(rr.Body).Decode(&act)
	Test.Equals(t, exp, act)

	// Test multiple portfolio values
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (1, '2018-10-11 21:05:29', 456)")
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (1, '2018-10-12 21:05:29', 789)")
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (2, '2018-10-10 21:05:29', 123)")
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (2, '2018-10-12 21:05:29', 456)")
	portfolioHandler.DB.Exec("INSERT INTO portfolio_history (user_id, date, value) VALUES (2, '2018-10-11 21:05:29', 789)")

	rr = Test.HandleRequest(req, portfolioHandler.GetPortfolioHistory)

	Test.Equals(t, http.StatusOK, rr.Code)

	date1, _ := time.Parse("2006-01-02T15:04:05Z", "2018-10-10T21:05:29Z")
	date2, _ := time.Parse("2006-01-02T15:04:05Z", "2018-10-11T21:05:29Z")
	date3, _ := time.Parse("2006-01-02T15:04:05Z", "2018-10-12T21:05:29Z")
	exp = portfolios.PortfolioHistory{
		UserID: 1,
		History: []portfolios.PortfolioValue{
			{ID: 1, Date: date1, Value: 123},
			{ID: 2, Date: date2, Value: 456},
			{ID: 3, Date: date3, Value: 789},
		},
	}
	act = portfolios.PortfolioHistory{}
	json.NewDecoder(rr.Body).Decode(&act)
	Test.Equals(t, exp, act)
}
