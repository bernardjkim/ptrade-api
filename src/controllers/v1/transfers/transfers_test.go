package transfers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bernardjkim/ptrade-api/pkg/types/transfers"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/transfers"
	"github.com/gorilla/mux"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	transferHandler TransferHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	transferHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = transferHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = transferHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = transferHandler.DB.Exec("ALTER TABLE orders AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = transferHandler.DB.Exec("ALTER TABLE transfer_orders AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestGetTransfersEmptyTable will test get transfers with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/transfers", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, transferHandler.GetTransfers)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

// TestGetTransfers
func TestGetTransfers(t *testing.T) {
	testSetup()

	// Test transfers initialization
	transferHandler.DB.Exec("INSERT INTO users (first, last, email, password) VALUES ('test1','test1','test1','test1')")
	transferHandler.DB.Exec("CALL new_transfer_order(1, 100)")

	req := httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/transfers", nil)

	// set mux vars
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)
	rr := Test.HandleRequest(req, transferHandler.GetTransfers)

	Test.Equals(t, http.StatusOK, rr.Code)

	transfer := transfers.TransferOrder{OrderID: 1, Balance: 100, Status: "FULFILLED"}

	exp := transfers.TransferOrders{UserID: 1, Transfers: []transfers.TransferOrder{transfer}}
	act := transfers.TransferOrders{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, exp.UserID, act.UserID)
	for index, _ := range exp.Transfers {
		Test.Equals(t, exp.Transfers[index].OrderID, act.Transfers[index].OrderID)
		Test.Equals(t, exp.Transfers[index].Balance, act.Transfers[index].Balance)
		Test.Equals(t, exp.Transfers[index].Status, act.Transfers[index].Status)
	}
}
