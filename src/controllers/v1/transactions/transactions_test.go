package transactions_test

import (
	"net/http"
	"strings"
	"testing"

	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"

	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/sessions"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/transactions"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/users"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	sessionHandler     SessionHandler
	transactionHandler TransactionHandler
	userHandler        UserHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	sessionHandler.Init(db)
	transactionHandler.Init(db)
	userHandler.Init(db)
}

// testSetup will clear user table and create a user to test with.
func testSetup() {
	Test.ClearTable("users")
	Test.ClearTable("transactions")
	createTestUser()
}

// createTestUser will create a user in the db with the values:
// 	{ID: 1, First: "John", Last: "Doe", Email: "johndoe@email.com", Password: "password"}
func createTestUser() {
	u := Test.User{ID: 1, First: "John", Last: "Doe", Email: "johndoe@email.com", Password: "password"}
	req := u.CreateUserReq()
	_ = Test.HandleRequest(req, userHandler.CreateUser)
}

// TestCreateTxn will test a simple create transaction request
func TestCreateTxn(t *testing.T) {
	testSetup()

	txn := Test.Transaction{UserID: 1, Quantity: 5}
	req := txn.CreateTxnReq(1, "AAPL")
	rr := Test.HandleRequest(req, transactionHandler.CreateTransaction)

	Test.Equals(t, http.StatusCreated, rr.Code)
}

// TestInvalidSymbolCreateTxn will test create transaction handler with an
// invalid symbol.
func TestInvalidSymbolCreateTxn(t *testing.T) {
	testSetup()

	txn := Test.Transaction{UserID: 1, Quantity: 5}
	req := txn.CreateTxnReq(1, "INVALID")
	rr := Test.HandleRequest(req, transactionHandler.CreateTransaction)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Unknown stock symbol."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

// TestUnauthorizedCreateTxn will test create transaction handler when
// a user attempts to create a transaction for another user.
func TestUnauthorizedCreateTxn(t *testing.T) {
	testSetup()

	txn := Test.Transaction{UserID: 1, Quantity: 5}
	req := txn.CreateTxnReq(2, "AAPL")
	rr := Test.HandleRequest(req, transactionHandler.CreateTransaction)

	Test.Equals(t, http.StatusUnauthorized, rr.Code)

	exp := "Unauthorized to make this request."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}
