package stocktransactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	StockTransactions "github.com/bernardjkim/ptrade-api/pkg/types/stock_transactions"
	Stocks "github.com/bernardjkim/ptrade-api/pkg/types/stocks"
	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/src/system/db"
)

// StockTransaction holds a transaction and stock object
type StockTransaction struct {
	Stock       Stocks.Stock                  `xorm:"extends" json:"stock"`
	Transaction StockTransactions.Transaction `xorm:"extends" json:"transaction"`
}

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

	var transactionList []StockTransaction

	//TODO: get table names from function?
	// get all transactions made by user join stock info
	t.DB.Table("stocks").Alias("s").
		Join("INNER", []string{"stock_transactions", "t"}, "s.id = t.stock_id").
		Where("t.user_id=?", userID).Find(&transactionList)

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

	// TODO: buy/sell depending on sign of quantity
	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid quantity value.", http.StatusBadRequest)
		return
	}

	symbol := r.FormValue("symbol")
	if len(symbol) < 1 {
		log.Println("Symbol not provided by user")
		http.Error(w, "No symbol provided.", http.StatusBadRequest)
		return
	}

	stock := Stocks.Stock{Symbol: symbol}
	if err = ORM.FindBy(t.DB, &stock); err != nil {
		log.Println(err)
		http.Error(w, "Unable to find stock in database.", http.StatusNotFound)
		return
	}

	if stock.ID < 1 {
		log.Printf("Unknown stock symbol: %s", symbol)
		http.Error(w, "Unknown stock symbol.", http.StatusBadRequest)
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
		http.Error(w, "Unable to retrieve stock price.", http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading body.", http.StatusInternalServerError)
		return
	}

	// TODO: format price to two dicimal places?
	price, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting price.", http.StatusInternalServerError)
		return
	}

	transaction := StockTransactions.Transaction{
		UserID:   userID,
		StockID:  stock.ID,
		Date:     timeStamp,
		Price:    price,
		Quantity: quantity,
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
