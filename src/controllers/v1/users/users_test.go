package users_test

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/users"
)

var (
	userHandler UserHandler
)

// init will initialize the request handlers needed for these test cases.
func init() {
	db := Test.InitTestDB()
	userHandler.Init(db)
}

// testSetup will run initial setup for each test case
func testSetup() {
	var err error
	_, err = userHandler.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = userHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal(err.Error())
	}
}

// TestEmptyTableGetUser will test the GetUser handler for an empty table.
// Should return status 200, but the body should not contain any user
// information besides the user id that was provided.
func TestEmptyTableGetUser(t *testing.T) {
	testSetup()

	u := Test.User{ID: 1}
	req := u.GetUserReq()
	rr := Test.HandleRequest(req, userHandler.GetUser)

	Test.Equals(t, http.StatusOK, rr.Code)

	act := Test.User{}
	json.NewDecoder(rr.Body).Decode(&act)

	Test.Equals(t, u, act)
}

// TestCreateNewUser will just test whether we can successfully create a new
// user in an empty table.
func TestCreateNewUser(t *testing.T) {
	testSetup()

	u := Test.User{First: "John", Last: "Doe", Email: "johndoe@email.com", Password: "password"}
	req := u.CreateUserReq()
	rr := Test.HandleRequest(req, userHandler.CreateUser)

	Test.Equals(t, http.StatusCreated, rr.Code)

	exp := Test.User{ID: 1, First: "John", Last: "Doe", Email: "johndoe@email.com"}
	act := Test.User{}
	json.NewDecoder(rr.Body).Decode(&act)
	Test.Equals(t, exp, act)
}

// TestCreateRepeatedUser will test the reponse when two requests to create
// a new user are made with the same email.
func TestCreateRepeatedUser(t *testing.T) {
	testSetup()

	u := Test.User{First: "John", Last: "Doe", Email: "johndoe@email.com", Password: "password"}
	req := u.CreateUserReq()
	_ = Test.HandleRequest(req, userHandler.CreateUser)
	rr := Test.HandleRequest(req, userHandler.CreateUser)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Email is already in use"
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}

// TestMissingFieldsCreateUser will test the reponse when email or password
// are missing from the request.
func TestMissingFieldsCreateUser(t *testing.T) {
	testSetup()

	// **Test missing email**

	missingEmail := Test.User{First: "John", Last: "Doe", Password: "password"}
	req := missingEmail.CreateUserReq()
	rr := Test.HandleRequest(req, userHandler.CreateUser)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Email and password are required."
	act := Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)

	// **Test missing password**

	missingPass := Test.User{First: "John", Last: "Doe", Email: "johndoe@email.com"}
	req = missingPass.CreateUserReq()
	rr = Test.HandleRequest(req, userHandler.CreateUser)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp = "Email and password are required."
	act = Test.ParseBody(rr.Body)
	Test.Equals(t, exp, act)
}
