package db

import (
	"time"
)

// UserTable represents a user stored in the database
type UserTable struct {
	ID       int64  `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	First    string `xorm:"VARCHAR(100) NOT NULL 'first'" json:"first" schema:"first"`
	Last     string `xorm:"VARCHAR(100) NOT NULL 'last'" json:"last" schema:"last"`
	Email    string `xorm:"VARCHAR(255) NOT NULL 'email'" json:"email" schema:"email"`
	Password string `xorm:"VARCHAR(255) NOT NULL 'password'" json:"password" schema:"password"`
}

// PortfolioHistoryTable reprents a snapshot of a users portfolio value at a given date
type PortfolioHistoryTable struct {
	ID     int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID int64     `xorm:"INTEGER 'user_id'" json:"user_id" schema:"user_id"`
	Date   time.Time `xorm:"DATETIME NOT NULL 'date'" json:"date" schema:"date"`
	Value  float64   `xorm:"FLOAT 'value'" json:"value" schema:"value"`
}

// StockTable represents a stock that is available in the database
type StockTable struct {
	ID            int64   `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	Symbol        string  `xorm:"VARCHAR(50) NOT NULL 'symbol'" json:"symbol" schema:"symbol"`
	Name          string  `xorm:"VARCHAR(255) NOT NULL 'name'" json:"name" schema:"name"`
	PricePerShare float64 `xorm:"FLOAT NOT NULL 'price_per_share'" json:"price_per_share" schema:"price_per_share"`
}

// OrderTable represents an order made by a user
type OrderTable struct {
	ID        int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID    int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	DateStart time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd   time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	Status    string    `xorm:"VARCHAR(20) NOT NULL 'status'" json:"status" schema:"status"`
}

// TransferOrderTable represents a transfer made by a user
type TransferOrderTable struct {
	ID      int64   `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	OrderID int64   `xorm:"INTEGER NOT NULL 'order_id'" json:"order_id" schema:"order_id"`
	Balance float64 `xorm:"FLOAT NOT NULL 'balance'" json:"balance" schema:"balance"`
}

// TradeOrderTable represents a trade made by a user
type TradeOrderTable struct {
	ID            int64   `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	OrderID       int64   `xorm:"INTEGER NOT NULL 'order_id'" json:"order_id" schema:"order_id"`
	StockID       int64   `xorm:"INTEGER NOT NULL 'stock_id'" json:"stock_id" schema:"stock_id"`
	Shares        int64   `xorm:"INTEGER NOT NULL 'shares'" json:"shares" schema:"shares"`
	PricePerShare float64 `xorm:"FLOAT NOT NULL 'price_per_share'" json:"price_per_share" schema:"price_per_share"`
}

// BalanceTable represents a users account balance
type BalanceTable struct {
	ID        int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID    int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	DateStart time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd   time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	Balance   float64   `xorm:"FLOAT NOT NULL 'balance'" json:"balance" schema:"balance"`
}

// PositionTable represents a users position for a specific stock
type PositionTable struct {
	ID        int64     `xorm:"SERIAL PRIMARY KEY 'id'" json:"id" schema:"id"`
	UserID    int64     `xorm:"INTEGER NOT NULL 'user_id'" json:"user_id" schema:"user_id"`
	DateStart time.Time `xorm:"DATETIME NOT NULL 'date_start'" json:"date_start" schema:"date_start"`
	DateEnd   time.Time `xorm:"DATETIME 'date_end'" json:"date_end" schema:"date_end"`
	StockID   int64     `xorm:"INTEGER NOT NULL 'stock_id'" json:"stock_id" schema:"stock_id"`
	Shares    int64     `xorm:"INTEGER NOT NULL 'shares'" json:"shares" schema:"shares"`
}

// TODO: xorm doesn't seem to support foriegn keys.
