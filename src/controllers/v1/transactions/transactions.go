package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	Stocks "github.com/bkim0128/bjstock-rest-service/pkg/types/stocks"
	Transactions "github.com/bkim0128/bjstock-rest-service/pkg/types/transactions"
	Users "github.com/bkim0128/bjstock-rest-service/pkg/types/users"
	ORM "github.com/bkim0128/bjstock-rest-service/src/system/db"

	"github.com/go-xorm/xorm"
	mux "github.com/gorilla/mux"
)

var db *xorm.Engine

// Init function will initialize this handler's connection to the db
func Init(DB *xorm.Engine) {
	db = DB
}

// GetTransactions returns an array of transactions made by a user
func GetTransactions(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id", http.StatusBadRequest)
		return
	}

	// if curUserID := r.Context().Value(Users.UserIDKey); curUserID != userID {
	// 	log.Println("Attepted to get another user's transactions")
	// 	http.Error(w, "Unauthorized to make this request", http.StatusBadRequest)
	// 	return
	// }

	// get all transactions made by user
	transactionList := []Transactions.Transaction{}

	// find all transactions with given user id
	if err := ORM.Find(db, &Transactions.Transaction{UserID: userID}, &transactionList); err != nil {
		log.Println(err)
		http.Error(w, "Unable to get transactions from database", http.StatusInternalServerError)
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
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id", http.StatusBadRequest)
		return
	}

	// get id of authenticated user
	if curUserID := r.Context().Value(Users.UserIDKey); curUserID != userID {
		log.Println("Attepted to create a transaction for another user")
		http.Error(w, "Unauthorized to make this request", http.StatusBadRequest)
		return
	}

	// TODO: buy/sell depending on sign of quantity
	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid quantity value", http.StatusBadRequest)
		return
	}

	symbol := r.FormValue("symbol")
	if len(symbol) < 1 {
		log.Println("Symbol not provided by user")
		http.Error(w, "No symbol provided", http.StatusBadRequest)
		return
	}

	stock := Stocks.Stock{Symbol: symbol}
	if err = ORM.FindBy(db, &stock); err != nil || userID < 1 {
		log.Println(err)
		http.Error(w, "Stock symbol not found", http.StatusNotFound)
		return
	}

	// TODO:
	// - what time and value is stored?
	// - clean up

	timeStamp := time.Now()
	fmt.Println("Time Stamp: ", timeStamp)

	// get current price for a share
	resp, err := http.Get("https://api.iextrading.com/1.0/stock/" + symbol + "/price")
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to retrieve stock price", http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	// TODO: format price to two dicimal places?
	price, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting price", http.StatusInternalServerError)
		return
	}

	transaction := Transactions.Transaction{
		UserID:   userID,
		StockID:  stock.ID,
		Date:     timeStamp,
		Price:    price,
		Quantity: quantity,
	}

	// store new transaction into database
	if err = ORM.Store(db, &transaction); err != nil {
		log.Println(err)
		http.Error(w, "Unable to make transaction", http.StatusInternalServerError)
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
