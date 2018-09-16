package sessions_test

import (
	"net/http"
	"strings"
	"testing"

	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/sessions"
	Test "github.com/bernardjkim/ptrade-api/src/controllers/v1/test"
	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/users"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var (
	userHandler    UserHandler
	sessionHandler SessionHandler
)

// init will initialize connection to database, and return a pointer to
// the xorm.Engine.
func init() {
	db := Test.InitTestDB()

	// initialize user handler
	userHandler.Init(db)
	sessionHandler.Init(db)
}

// testSetup will clear user table and create a user to test with.
func testSetup() {
	Test.ClearUserTable()
	createTestUser()

}

// createTestUser will create a user in the db with the values:
// 	{ID: 1, First: "John", Last: "Doe", Email: "johndoe@email.com", Password: "password"}
func createTestUser() {
	u := Test.User{1, "John", "Doe", "johndoe@email.com", "password"}
	req := u.CreateUserReq()
	_ = Test.HandleRequest(req, userHandler.CreateUser)
}

func TestCreateSession(t *testing.T) {
	testSetup()

	req := Test.CreateLoginReq("johndoe@email.com", "password")
	rr := Test.HandleRequest(req, sessionHandler.CreateSession)

	Test.Equals(t, http.StatusCreated, rr.Code)

	// TODO: how to test body response for jwt token?
	// for now just testing status code.
	// exp := "jwt token"
	// act := strings.TrimSuffix(rr.Body.String(), "\n")
	// equals(t, exp, act)
}

func TestInvalidEmailCreateSession(t *testing.T) {
	testSetup()

	req := Test.CreateLoginReq("invalid@email.com", "password")
	rr := Test.HandleRequest(req, sessionHandler.CreateSession)

	Test.Equals(t, http.StatusNotFound, rr.Code)

	exp := "No user with provided email exists."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

func TestInvalidPassCreateSession(t *testing.T) {
	testSetup()

	req := Test.CreateLoginReq("johndoe@email.com", "invalid")
	rr := Test.HandleRequest(req, sessionHandler.CreateSession)

	Test.Equals(t, http.StatusUnauthorized, rr.Code)

	exp := "Credentials do not match."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}

func TestMissingFieldsCreateSession(t *testing.T) {
	testSetup()

	missingEmail := Test.CreateLoginReq("", "password")
	rr := Test.HandleRequest(missingEmail, sessionHandler.CreateSession)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp := "Email and password are required."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)

	missingPass := Test.CreateLoginReq("johndoe@email.com", "")
	rr = Test.HandleRequest(missingPass, sessionHandler.CreateSession)

	Test.Equals(t, http.StatusBadRequest, rr.Code)

	exp = "Email and password are required."
	act = strings.TrimSuffix(rr.Body.String(), "\n")
	Test.Equals(t, exp, act)
}
