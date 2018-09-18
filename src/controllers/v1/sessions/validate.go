package sessions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bernardjkim/ptrade-api/src/system/jwt"
)

// Validate users through session tokens. Respond to client with user data.
func (s *SessionHandler) Validate(w http.ResponseWriter, r *http.Request) {

	// check if token is present
	tokenVal := r.Header.Get("Session-Token")
	if len(tokenVal) < 1 {
		log.Println("Ignoring request. No token present.")
		http.Error(w, "Token required for check.", http.StatusUnauthorized)
		return
	}

	user, err := jwt.GetUserFromToken(s.DB, tokenVal)
	if err != nil {
		log.Println(err)
		// TODO: is it possible to get an error when a valid token is given??
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "Invalid Token.", http.StatusBadRequest)
		return
	}

	user.Password = ""

	// login := LoginData{User: user, Token: tokenVal}

	// convert packet to JSON
	packet, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
