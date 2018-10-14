package balances

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	"github.com/bernardjkim/ptrade-api/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/src/system/db"
)

// BalanceHandler struct needs to be initialized with a database connection.
type BalanceHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (h *BalanceHandler) Init(DB *xorm.Engine) {
	h.DB = DB
}

// GetBalance returns the current user's account balance
func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
		return
	}

	// check if user exists
	exists, err := ORM.Exists(h.DB, &users.User{ID: userID})
	if err != nil || !exists {
		log.Println("Provided user does not exist")
		http.Error(w, "Provided user id does not exist in databse", http.StatusBadRequest)
		return
	}

	bal, err := ORM.GetBalance(h.DB, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to get balance", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(bal)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}

// // CreateTransaction will execute user transactions and update the database with the
// // desired transaction.
// func (t *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	// get user id from url
// 	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
// 		return
// 	}

// 	// get id of authenticated user
// 	if curUserID := r.Context().Value(Users.UserIDKey); curUserID != userID {
// 		log.Printf("Attepted to create a txn for user: %d, authenticated as user: %d\n", userID, curUserID)
// 		http.Error(w, "Unauthorized to make this request.", http.StatusUnauthorized)
// 		return
// 	}

// 	value, err := strconv.ParseFloat(r.FormValue("value"), 64)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "No value provided.", http.StatusBadRequest)
// 		return
// 	}

// 	timeStamp := time.Now()

// 	transaction := BankingTransactions.Transaction{
// 		UserID: userID,
// 		Date:   timeStamp,
// 		Value:  value,
// 	}

// 	// store new transaction into database
// 	if err = ORM.Store(t.DB, &transaction); err != nil {
// 		log.Println(err)
// 		http.Error(w, "Unable to make transaction.", http.StatusInternalServerError)
// 		return
// 	}

// 	// convert packet to JSON
// 	packet, err := json.Marshal(transaction)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(packet)
// }
