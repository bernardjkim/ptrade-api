package users_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	. "github.com/bernardjkim/ptrade-api/src/controllers/v1/users"
	DB "github.com/bernardjkim/ptrade-api/src/system/db"
)

// NOTE: trimming reponse body of \n because http.Error calls Fprintln which
// adds a new line to the end of the error msg.

var userHandler UserHandler

// init will initialize connection to database, and return a pointer to
// the xorm.Engine.
func init() {
	_ = godotenv.Load(os.Getenv("GOPATH") + "/src/github.com/bernardjkim/ptrade-api/.env")

	url := os.Getenv("JAWSDB_ORANGE_URL")
	if len(url) < 0 {
		log.Println("Unable to get database url")
	}

	db, err := DB.ConnectURL(url, "parseTime=true")
	if err != nil {
		log.Println("Unable to connect to db")
		panic(err)
	}

	DB.Init(db)

	// initialize user handler
	userHandler.Init(db)
}

// clearUserTable will delete all entries in the 'users' table and reset the
// auto increment to start at 1.
func clearUserTable() {
	userHandler.DB.Exec("DELETE FROM users")
	userHandler.DB.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

// handleRequest will handle the given request with the given handlerFunc and
// return a *httptest.ResponseRecorder.
func handleRequest(req *http.Request, handlerFunc http.HandlerFunc) (rr *httptest.ResponseRecorder) {
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)
	return
}

// **TESTING USER HANDLER FUNCTIONS**

// TestEmptyTableGetUser will test the GetUser handler for an empty table.
// Should return status 200, but the body should not contain any user
// information besides the user id that was provided.
func TestEmptyTableGetUser(t *testing.T) {
	clearUserTable()

	// get user of id 1
	u := user{1, "", "", "", ""}
	req := u.getUserReq()
	rr := handleRequest(req, userHandler.GetUser)

	// Verify response code
	equals(t, http.StatusOK, rr.Code)

	// Verify response body
	exp := u.toString()
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	equals(t, exp, act)
}

// TestCreateNewUser will just test whether we can successfully create a new
// user in an empty table.
func TestCreateNewUser(t *testing.T) {
	clearUserTable()

	// **Test missing email**
	u := user{1, "John", "Doe", "johndoe@email.com", "password"}
	req := u.createUserReq()
	rr := handleRequest(req, userHandler.CreateUser)

	// Verify response code
	equals(t, http.StatusCreated, rr.Code)

	// Verify response body
	u.Password = "" // password will be dropped in response
	exp := u.toString()
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	equals(t, exp, act)
}

// TestCreateRepeatedUser will test the reponse when two requests to create
// a new user are made with the same email.
func TestCreateRepeatedUser(t *testing.T) {
	clearUserTable()

	u := user{1, "John", "Doe", "johndoe@email.com", "password"}
	req := u.createUserReq()
	_ = handleRequest(req, userHandler.CreateUser)
	rr := handleRequest(req, userHandler.CreateUser)

	// Verify response code
	equals(t, http.StatusBadRequest, rr.Code)

	// Verify response body
	exp := "Email is already in use"
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	equals(t, exp, act)
}

// TestMissingFieldsCreateUser will test the reponse when email or password
// are missing from the request.
func TestMissingFieldsCreateUser(t *testing.T) {
	clearUserTable()

	// **Test missing email**
	u := user{0, "John", "Doe", "johndoe@email.com", "password"}
	req := u.createUserReq()
	req.Form.Set("email", "")
	rr := handleRequest(req, userHandler.CreateUser)

	// Verify response code
	equals(t, http.StatusBadRequest, rr.Code)

	// Verify response body
	exp := "Email and password are required."
	act := strings.TrimSuffix(rr.Body.String(), "\n")
	equals(t, exp, act)

	// **Test missing password**

	req.Form.Set("email", "johndoe@email.com")
	req.Form.Set("password", "")
	rr = handleRequest(req, userHandler.CreateUser)

	// Verify response code
	equals(t, http.StatusBadRequest, rr.Code)

	// Verify response body
	exp = "Email and password are required."
	act = strings.TrimSuffix(rr.Body.String(), "\n")
	equals(t, exp, act)
}

type user struct {
	ID       int64
	First    string
	Last     string
	Email    string
	Password string
}

// createUserReq will return a *http.Request to serve to the CreateUser handler
func (u *user) createUserReq() (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/users", nil)
	req.Form = url.Values{
		"first":    {u.First},
		"last":     {u.Last},
		"email":    {u.Email},
		"password": {u.Password},
	}
	return
}

// getUserReq will return a *http.Request to serve to the GetUser handler
func (u *user) getUserReq() (req *http.Request) {
	req = httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}", nil)
	vars := map[string]string{
		"ID": strconv.Itoa(int(u.ID)),
	}
	req = mux.SetURLVars(req, vars)
	return
}

func (u *user) toString() (s string) {
	s = `{"id":` + strconv.Itoa(int(u.ID)) +
		`,"first":"` + u.First +
		`","last":"` + u.Last +
		`","email":"` + u.Email +
		`","password":"` + u.Password + `"}`
	return
}

//
// ** https://github.com/benbjohnson/testing **
//

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		fmt.Println(exp)
		fmt.Println(act)
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
