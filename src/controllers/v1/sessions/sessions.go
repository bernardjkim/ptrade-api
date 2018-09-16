package sessions

import (
	"github.com/go-xorm/xorm"
)

// LoginData holds user token and data
// type LoginData struct {
// 	Token string     `json:"token"`
// 	User  Users.User `json:"user"`
// }

// SessionHandler struct needs to be initialized with a database connection.
type SessionHandler struct {
	DB *xorm.Engine
}

// Init function initializes this sessions db connection
func (s *SessionHandler) Init(DB *xorm.Engine) {
	s.DB = DB
}
