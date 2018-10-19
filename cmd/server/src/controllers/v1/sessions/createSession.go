package sessions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bernardjkim/ptrade-api/cmd/server/src/system/jwt"

	Users "github.com/bernardjkim/ptrade-api/cmd/server/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/cmd/server/src/system/db"
	Passwords "github.com/bernardjkim/ptrade-api/cmd/server/src/system/passwords"
)

//TODO: status codes

// CreateSession function attempts to authenticate user with given credentials.
// Responds with jwt token if successful.
func (s *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")

	// check if both email and password arguments have been provided
	if len(email) < 1 || len(password) < 1 {
		http.Error(w, "Email and password are required.", http.StatusBadRequest)
		return
	}

	user := Users.User{Email: email}
	if err := ORM.FindBy(s.DB, &user); err != nil {
		log.Println(err)
		http.Error(w, "Unable to find users.", http.StatusInternalServerError)
		return
	}

	// check if user with matching email exists in database
	if user.ID < 1 {
		log.Println(w, "Invalid email")
		http.Error(w, "No user with provided email exists.", http.StatusNotFound)
		return
	}

	// check for valid password
	if !Passwords.IsValid(user.Password, password) {
		log.Println(w, "Invalid password")
		http.Error(w, "Credentials do not match.", http.StatusUnauthorized)
		return
	}

	// generate authentication token for user
	token := jwt.GetToken(user.ID)

	// set cookie header
	http.SetCookie(w, &http.Cookie{
		Name:       "api.example.com",
		Value:      token,
		Path:       "/",
		RawExpires: "3600", // one hour expiration time
	})

	// convert packet to JSON
	packet, err := json.Marshal(struct {
		SessionToken string `json:"Session-Token"`
	}{token})
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}
