package trades

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	Stocks "github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/stocks"
	Users "github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/cmd/server/src/system/db"
)

// TradeHandler struct needs to be initialized with a database connection.
type TradeHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (h *TradeHandler) Init(DB *xorm.Engine) {
	h.DB = DB
}

// CreateTrade creates a new trade order
func (h *TradeHandler) CreateTrade(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	symbol := r.FormValue("symbol")
	shares := r.FormValue("shares")

	// verify that all fields have been provided
	if len(symbol) < 1 || len(shares) < 1 {
		http.Error(w, "Symbol and number of shares are required.", http.StatusBadRequest)
		return
	}

	stock := Stocks.Stock{Symbol: symbol}

	// check if symbol exists
	err := ORM.FindBy(h.DB, &stock)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if stock.ID < 1 {
		log.Println("Provided symbol id does not exist in database")
		http.Error(w, "Provided symbol id does not exist in database", http.StatusBadRequest)
		return

	}

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
		return
	}

	// get id of authenticated user
	curUserID := r.Context().Value(Users.UserIDKey)
	if curUserID != userID {
		log.Printf("Attempted to create order for user: %d, authenticated as user: %d\n", userID, curUserID)
		http.Error(w, "Unauthorized to make this request.", http.StatusUnauthorized)
		return
	}

	// check if user exists
	userExists, err := ORM.Exists(h.DB, &Users.User{ID: userID})
	if err != nil || !userExists {
		log.Println("Provided user does not exist")
		http.Error(w, "Provided user id does not exist in databse", http.StatusBadRequest)
		return
	}

	numShares, err := strconv.ParseInt(shares, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to parse int", http.StatusInternalServerError)
		return
	}

	err = ORM.SetPricePerShare(h.DB, symbol)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err = ORM.NewTrade(h.DB, userID, stock.ID, numShares); err != nil {
		log.Println(err)
		http.Error(w, "Unable to create new trade", http.StatusInternalServerError)
		return
	}

	// trades, err := ORM.GetTrades(h.DB, userID)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Unable to get trades", http.StatusInternalServerError)
	// 	return
	// }

	// // convert packet to JSON
	// packet, err := json.Marshal(trades)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	// w.Write(packet)
}

// GetTrades returns the current user's postions
func (h *TradeHandler) GetTrades(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id.", http.StatusBadRequest)
		return
	}

	// check if user exists
	exists, err := ORM.Exists(h.DB, &Users.User{ID: userID})
	if err != nil || !exists {
		log.Println("Provided user does not exist")
		http.Error(w, "Provided user id does not exist in databse", http.StatusBadRequest)
		return
	}

	trades, err := ORM.GetTrades(h.DB, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to get trades", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(trades)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
