package transfers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bernardjkim/ptrade-api/pkg/types/transfers"
	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/transfers"
	"github.com/gorilla/mux"
)

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

// getTransferReq returns a new transfer req with the given fields
func getTransfersReq(id string) (req *http.Request) {
	req = httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}/transfers", nil)
	// set mux vars
	vars := map[string]string{
		"ID": id,
	}
	req = mux.SetURLVars(req, vars)
	return
}

// createTransferReq returns a new transfer req with the given fields
func createTransferReq(id, balance string) (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/users/{ID:[0-9]+}/transfers", nil)

	req.Form = url.Values{
		"balance": {balance},
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

// transferOrder returns a new TransferOrder with the given fields
func transferOrder(orderID int64, balance float64, status string) (transfer transfers.TransferOrder) {
	transfer = transfers.TransferOrder{OrderID: orderID, Balance: balance, Status: status}
	return
}

// compareTransfers will compare the fields of the exp and act
func compareTransfers(t *testing.T, exp transfers.TransferOrders, act transfers.TransferOrders) {
	Test.Equals(t, fmt.Sprintf("%s: %d", "user id", exp.UserID), fmt.Sprintf("%s: %d", "user id", act.UserID))
	for index, _ := range exp.Transfers {
		Test.Equals(t, fmt.Sprintf("%s: %d", "order id", exp.Transfers[index].OrderID), fmt.Sprintf("%s: %d", "order id", act.Transfers[index].OrderID))
		Test.Equals(t, fmt.Sprintf("%s: %.2f", "balance", exp.Transfers[index].Balance), fmt.Sprintf("%s: %.2f", "balance", act.Transfers[index].Balance))
		Test.Equals(t, fmt.Sprintf("%s: %s", "status", exp.Transfers[index].Status), fmt.Sprintf("%s: %s", "status", act.Transfers[index].Status))
	}
}

// TestGetTransfersEmptyTable will test get transfers with an empty table
func TestGetHistoryEmptyTable(t *testing.T) {
	testSetup()

	req := getTransfersReq("1")
	rr := Test.HandleRequest(req, transferHandler.GetTransfers)
	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Provided user id does not exist in databse"
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}

// TestGetTransfers
func TestGetTransfers(t *testing.T) {
	testSetup()

	// Test transfers initialization
	Test.NewUser("test1", "test1", "test1", "test1")
	req := createTransferReq("1", "100.00")
	_ = Test.HandleRequest(req, transferHandler.CreateTransfer)

	req = getTransfersReq("1")
	rr := Test.HandleRequest(req, transferHandler.GetTransfers)
	Test.Equals(t, http.StatusOK, rr.Code)

	transferList := []transfers.TransferOrder{}
	transferList = append(transferList, transferOrder(1, 100.00, "FULFILLED"))

	exp := transfers.TransferOrders{UserID: 1, Transfers: transferList}
	act := transfers.TransferOrders{}
	json.NewDecoder(rr.Body).Decode(&act)
	compareTransfers(t, exp, act)
}

// TestCreateTransfer will test creating transfer orders
func TestCreateTransfer(t *testing.T) {
	testSetup()
	Test.NewUser("test1", "test1", "test1", "test1")

	// Test single entry
	req := createTransferReq("1", "100.00")
	rr := Test.HandleRequest(req, transferHandler.CreateTransfer)
	Test.Equals(t, http.StatusCreated, rr.Code)

	transferList := []transfers.TransferOrder{}
	transferList = append(transferList, transferOrder(1, 100.00, "FULFILLED"))

	exp := transfers.TransferOrders{UserID: 1, Transfers: transferList}
	act := transfers.TransferOrders{}
	json.NewDecoder(rr.Body).Decode(&act)
	compareTransfers(t, exp, act)

	// Test multiple entries
	req = createTransferReq("1", "123")
	_ = Test.HandleRequest(req, transferHandler.CreateTransfer)
	req = createTransferReq("1", "456")
	_ = Test.HandleRequest(req, transferHandler.CreateTransfer)
	req = createTransferReq("1", "789.123")
	rr = Test.HandleRequest(req, transferHandler.CreateTransfer)
	Test.Equals(t, http.StatusCreated, rr.Code)

	transferList = append(transferList, transferOrder(2, 123, "FULFILLED"))
	transferList = append(transferList, transferOrder(3, 456, "FULFILLED"))
	transferList = append(transferList, transferOrder(4, 789.123, "FULFILLED"))

	exp = transfers.TransferOrders{UserID: 1, Transfers: transferList}
	act = transfers.TransferOrders{}
	json.NewDecoder(rr.Body).Decode(&act)
	compareTransfers(t, exp, act)
}
