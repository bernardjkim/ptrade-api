package portfolios

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bernardjkim/ptrade-api/pkg/types/users"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	ORM "github.com/bernardjkim/ptrade-api/src/system/db"
)

// PortfolioHandler struct needs to be initialized with a database connection.
type PortfolioHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (h *PortfolioHandler) Init(DB *xorm.Engine) {
	h.DB = DB
}

// GetPortfolioHistory returns a user's portfolio history
func (h *PortfolioHandler) GetPortfolioHistory(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id", http.StatusBadRequest)
		return
	}

	// check if user exists
	exists, err := ORM.Exists(h.DB, &users.User{ID: userID})
	if err != nil || !exists {
		log.Println("Provided user does not exist")
		http.Error(w, "Provided user id does not exist in databse", http.StatusBadRequest)
		return
	}

	history, err := ORM.GetPortfolioHistory(h.DB, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to get portfolio history", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(history)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
