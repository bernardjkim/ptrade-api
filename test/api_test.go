package api_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	UserHandler "github.com/bkim0128/bjstock-rest-service/src/controllers/v1/users"
	DB "github.com/bkim0128/bjstock-rest-service/src/system/db"
)

// initDB will initialize connection to database, and return a pointer to
// the xorm.Engine.
func initDB() (db *xorm.Engine) {
	_ = godotenv.Load(os.Getenv("GOPATH") + "/src/github.com/bkim0128/bjstock-rest-service/.env")

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
	return
}

// TestEmptyTableGetUser will test the GetUser handler for an empty table.
// Should return status 200, but the body should not contain any user
// information besides the user id that was provided.
func TestEmptyTableGetUser(t *testing.T) {

	db := initDB()
	UserHandler.Init(db)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/v1/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserHandler.GetUser)

	// set url variables
	vars := map[string]string{
		"ID": "1",
	}
	req = mux.SetURLVars(req, vars)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":1,"first":"","last":"","email":"","password":""}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
