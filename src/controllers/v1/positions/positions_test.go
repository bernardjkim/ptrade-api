package positions_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/types/positions"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/positions"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	"github.com/gorilla/mux"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	positionHandler PositionHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	positionHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = positionHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = positionHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = positionHandler.DB.Exec("ALTER TABLE positions AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestGetPositionsEmptyTable will test get position with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/positions", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, positionHandler.GetPositions)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

// TestGetPositions
func TestGetPositions(t *testing.T) {
	testSetup()

	// Test positions initialization
	positionHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")
	positionHandler.DB.Exec("INSERT INTO positions (user_id, stock_id, shares) VALUES (1, 1, 5)")

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/positions", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, positionHandler.GetPositions)

	Test.Equals(t, http.StatusOK, rr.Code)

	pos := positions.Position{StockID: 1, Date: time.Now(), Shares: 5}

	exp := positions.Positions{UserID: 1, Positions: []positions.Position{pos}}
	act := positions.Positions{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, exp.UserID, act.UserID)
	for index, _ := range exp.Positions {
		Test.Equals(t, exp.Positions[index].StockID, act.Positions[index].StockID)
		Test.Equals(t, exp.Positions[index].Shares, act.Positions[index].Shares)
	}
}
