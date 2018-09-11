package sessions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bkim0128/bjstock-rest-service/src/system/jwt"

	Users "github.com/bkim0128/bjstock-rest-service/pkg/types/users"
	ORM "github.com/bkim0128/bjstock-rest-service/src/system/db"
	Passwords "github.com/bkim0128/bjstock-rest-service/src/system/passwords"
)

//TODO: status codes

// CreateSession function attempts to authenticate user with given credentials.
// Responds with jwt token if successful.
func CreateSession(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")

	// check if both email and password arguments have been provided
	if len(email) < 1 || len(password) < 1 {
		http.Error(w, "Email and password are required.", http.StatusBadRequest)
		return
	}

	// check if user with matching email exists in database
	user := Users.User{Email: email}
	if err := ORM.FindBy(db, &user); err != nil || user.ID < 1 {
		log.Println(err)
		http.Error(w, "Credentials do not match.", http.StatusUnauthorized)
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
	packet, err := json.Marshal(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(packet)
}