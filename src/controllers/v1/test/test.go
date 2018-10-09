package test

import (
	"context"
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
	"testing"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	DB "github.com/bernardjkim/ptrade-api/src/system/db"
)

// User struct for testing user handler
type User struct {
	ID       int
	First    string
	Last     string
	Email    string
	Password string
}

// Transaction struct testing transaction handler
type Transaction struct {
	ID       int64
	UserID   int64
	StockID  int64
	Date     time.Time
	Price    float64
	Quantity int64
}

var testDB *xorm.Engine

// InitTestDB will initialize connection to database, and return a pointer to
// the xorm.Engine.
func InitTestDB() *xorm.Engine {
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

	// DB.Init(db)
	testDB = db
	return testDB
}

// HandleRequest will handle the given request with the given handlerFunc and
// return a *httptest.ResponseRecorder.
func HandleRequest(req *http.Request, handlerFunc http.HandlerFunc) (rr *httptest.ResponseRecorder) {
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, req)
	return
}

// ClearTable will delete all entries in the {tableName} table and reset the
// auto increment to start at 1.
func ClearTable(tableName string) {
	testDB.Exec("DELETE FROM ?", tableName)
	testDB.Exec("ALTER TABLE ? AUTO_INCREMENT = 1", tableName)
}

// ** USER REQUESTS **

// CreateUserReq will return a *http.Request to serve to the CreateUser handler
// func CreateUserReq(first, last, email, password string) (req *http.Request) {
func (u *User) CreateUserReq() (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/users", nil)
	req.Form = url.Values{
		"first":    {u.First},
		"last":     {u.Last},
		"email":    {u.Email},
		"password": {u.Password},
	}
	return
}

// GetUserReq will return a *http.Request to serve to the GetUser handler
func (u *User) GetUserReq() (req *http.Request) {
	req = httptest.NewRequest("GET", "/v1/users/{ID:[0-9]+}", nil)
	vars := map[string]string{
		"ID": strconv.Itoa(u.ID),
	}
	req = mux.SetURLVars(req, vars)
	return
}

// ** SESSION REQUESTS **

// CreateLoginReq will return a *http.Request to serve to the CreateSession handler
func CreateLoginReq(email, password string) (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/sessions", nil)
	req.Form = url.Values{
		"email":    {email},
		"password": {password},
	}
	return
}

// ** TRANSACTION REQUESTS **

// CreateTxnReq will return a *http.Request to serve to CreateTransaction handler
func (t *Transaction) CreateTxnReq(userID int64, symbol string) (req *http.Request) {
	req = httptest.NewRequest("POST", "/v1/users/{ID:[0-9]+}/transactions", nil)

	// set form values
	req.Form = url.Values{
		"quantity": {strconv.Itoa(int(t.Quantity))},
		"symbol":   {symbol},
	}

	// set mux vars
	vars := map[string]string{
		"ID": strconv.Itoa(int(t.UserID)),
	}
	req = mux.SetURLVars(req, vars)

	// set user id in context
	ctx := context.WithValue(req.Context(), Users.UserIDKey, userID)
	req = req.WithContext(ctx)

	return
}

// ** https://github.com/benbjohnson/testing **

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		fmt.Println(exp)
		fmt.Println(act)
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
