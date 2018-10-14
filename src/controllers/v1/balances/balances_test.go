package balances_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bernardjkim/ptrade-api/pkg/types/balances"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/balances"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	"github.com/gorilla/mux"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	balanceHandler BalanceHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	balanceHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = balanceHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = balanceHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = balanceHandler.DB.Exec("ALTER TABLE balances AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestGetBalanceEmptyTable will test get balance with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/balance", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, balanceHandler.GetBalance)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

// TestGetBalance
func TestGetBalance(t *testing.T) {
	testSetup()

	// Test balance initialization
	balanceHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")
	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/balance", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, balanceHandler.GetBalance)

	Test.Equals(t, http.StatusOK, rr.Code)

	exp := balances.Balance{ID: 1, UserID: 1, DateStart: time.Now(), DateEnd: time.Now(), Balance: 0}
	act := balances.Balance{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, exp.ID, act.ID)
	Test.Equals(t, exp.UserID, act.UserID)
	Test.Equals(t, exp.Balance, act.Balance)
}
