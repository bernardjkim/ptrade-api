package users

import (
	"github.com/bernardjkim/ptrade-api/cmd/server/pkg/db"
)

type Users []User

type User db.UserTable

type key string

const UserIDKey key = "userID"

func (u *User) TableName() string {
	return "users"
}
