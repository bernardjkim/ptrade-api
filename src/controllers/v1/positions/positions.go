package positions

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

// PositionHandler struct needs to be initialized with a database connection.
type PositionHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (h *PositionHandler) Init(DB *xorm.Engine) {
	h.DB = DB
}

// GetPositions returns the current user's postions
func (h *PositionHandler) GetPositions(w http.ResponseWriter, r *http.Request) {

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

	positions, err := ORM.GetPositions(h.DB, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to get positions", http.StatusInternalServerError)
		return
	}

	// convert packet to JSON
	packet, err := json.Marshal(positions)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
