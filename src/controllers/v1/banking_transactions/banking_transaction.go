package bankingtransactions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	BankingTransactions "github.com/bernardjkim/ptrade-api/pkg/types/banking_transactions"
	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/src/system/db"
)

// TransactionHandler struct needs to be initialized with a database connection.
type TransactionHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (t *TransactionHandler) Init(DB *xorm.Engine) {
	t.DB = DB
}

// GetTransactions returns an array of transactions made by a user
func (t *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
		return
	}

	var transactionList []BankingTransactions.Transaction

	// get list of available stocks from database
	if err := ORM.Find(t.DB, &BankingTransactions.Transaction{UserID: userID}, &transactionList); err != nil {
		log.Println(err)
		http.Error(w, "Unable to get banking transactions", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(transactionList)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}

// CreateTransaction will execute user transactions and update the database with the
// desired transaction.
func (t *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
		return
	}

	// get id of authenticated user
	if curUserID := r.Context().Value(Users.UserIDKey); curUserID != userID {
		log.Printf("Attepted to create a txn for user: %d, authenticated as user: %d\n", userID, curUserID)
		http.Error(w, "Unauthorized to make this request.", http.StatusUnauthorized)
		return
	}

	value, err := strconv.ParseFloat(r.FormValue("value"), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "No value provided.", http.StatusBadRequest)
		return
	}

	timeStamp := time.Now()

	transaction := BankingTransactions.Transaction{
		UserID: userID,
		Date:   timeStamp,
		Value:  value,
	}

	// store new transaction into database
	if err = ORM.Store(t.DB, &transaction); err != nil {
		log.Println(err)
		http.Error(w, "Unable to make transaction.", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(transaction)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
