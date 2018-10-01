package db

import (
	"time"
)

// UserTable represents a user stored in the database
type UserTable struct {
	ID       int64  `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	First    string `xorm:"VARCHAR(50) NOT NULL 'first'" json:"first" schema:"first"`
	Last     string `xorm:"VARCHAR(50) NOT NULL 'last'" json:"last" schema:"last"`
	Email    string `xorm:"VARCHAR(50) NOT NULL 'email'" json:"email" schema:"email"`
	Password string `xorm:"VARCHAR(255) NOT NULL 'password'" json:"password" schema:"password"`
}

// StockTable represents a stock that is available in the database
type StockTable struct {
	ID     int64  `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	Symbol string `xorm:"VARCHAR(50) NOT NULL 'symbol'" json:"symbol" schema:"symbol"`
	Name   string `xorm:"VARCHAR(255) NOT NULL 'name'" json:"name" schema:"name"`
}

// StockTransactionTable represents a stock transaction made by a user
type StockTransactionTable struct {
	ID       int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID   int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	StockID  int64     `xorm:"INTEGER NOT NULL 'stock_id'" json:"stock_id" schema:"stock_id"`
	Date     time.Time `xorm:"DATETIME NOT NULL 'date'" json:"date" schema:"date"`
	Price    float64   `xorm:"FLOAT NOT NULL 'price'" json:"price" schema:"price"`
	Quantity int64     `xorm:"INTEGER NOT NULL 'quantity'" json:"quantity" schema:"quantity"`
}

// BankingTransactionTable represents a banking transaction made by a user
type BankingTransactionTable struct {
	ID     int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	Date   time.Time `xorm:"DATETIME NOT NULL 'date'" json:"date" schema:"date"`
	Value  float64   `xorm:"FLOAT NOT NULL 'value'" json:"value" schema:"value"`
}

// TODO: xorm doesn't seem to support foriegn keys.
