package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"

	Users "github.com/bernardjkim/ptrade-api/pkg/types/users"
	ORM "github.com/bernardjkim/ptrade-api/src/system/db"
	Passwords "github.com/bernardjkim/ptrade-api/src/system/passwords"
)

// UserHandler struct needs to be initialized with a database connection.
type UserHandler struct {
	DB *xorm.Engine
}

// Init function will initialize this handler's connection to the db
func (h *UserHandler) Init(DB *xorm.Engine) {
	h.DB = DB
}

// GetUsers responds with a list of all users?
func GetUsers(w http.ResponseWriter, r *http.Request) {
}

// GetUser responds with user information
// TODO: return error if no user with give id? or just empty user?
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	// get user id from url
	userID, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "Provided invalid id", http.StatusBadRequest)
		return
	}

	user := Users.User{ID: userID}

	// find all transactions with given user id
	if err := ORM.FindBy(h.DB, &user); err != nil {
		log.Println(err)
		http.Error(w, "Unable to get user from database", http.StatusInternalServerError)
		return
	}

	// remove password field
	user.Password = ""

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

// TODO: which status code should be returned

// CreateUser function attempts to create a new user profile.
// Responds with status 201 OK if successfully created new user. Will respond
// with an error code otherwise.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	first := r.FormValue("first")
	last := r.FormValue("last")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// verify that email and password have been provided
	if len(email) < 1 || len(password) < 1 {
		http.Error(w, "Email and password are required.", http.StatusBadRequest)
		return
	}

	user := Users.User{Email: email}
	if err := ORM.FindBy(h.DB, &user); err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// verify that there does not already exist a user with the same email
	if user.ID > 0 {
		log.Println("Attempted to create new user with preexisting email")
		http.Error(w, "Email is already in use", http.StatusBadRequest)
		return
	}

	// encrypt password
	encryptedPassword, err := Passwords.Encrypt(password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to encrypt password", http.StatusInternalServerError)
		return
	}

	user.First = first
	user.Last = last
	user.Password = encryptedPassword

	// store new user into database
	if err = ORM.Store(h.DB, &user); err != nil {
		log.Println(err)
		http.Error(w, "Unable to create account", http.StatusInternalServerError)
		return
	}

	// remove password in response
	user.Password = ""

	packet, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(packet)
}
