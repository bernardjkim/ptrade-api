package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/joho/godotenv"

	DB "github.com/bernardjkim/ptrade-api/src/system/db"
)

var (
	dbURL     string
	dboptions string
)

func init() {
	flag.StringVar(&dboptions, "dboptions", "parseTime=true", "Set the port for the application")

	flag.Parse()

	_ = godotenv.Load()

	if url := os.Getenv("JAWSDB_URL"); len(url) > 0 {
		dbURL = url
	}
}

func main() {
	currentTime := getEST(time.Now())

	// stock market closed on weekends
	if currentTime.Weekday().String() == "Saturday" ||
		currentTime.Weekday().String() == "Sunday" {
		return
	}

	// stock market hours from 6:00am - 3:00 pm ??
	if currentTime.Hour() < 6 || currentTime.Hour() >= 15 {
		return
	}

	db, err := DB.ConnectURL(dbURL, dboptions)
	if err != nil {
		log.Println("Unable to connect to db")
		panic(err)
	}
	db.ShowSQL()

	symbols, err := getSymbols(db)
	if err != nil {
		log.Println("Unalge to get symbols")
		panic(err)
	}

	getPrices(db, symbols)

}

// getEST converts t to EST
func getEST(t time.Time) time.Time {
	location, err := time.LoadLocation("EST")
	if err != nil {
		log.Println(err)
	}
	return t.In(location)
}

func getSymbols(DB *xorm.Engine) (symbols []string, err error) {
	rows, err := DB.QueryString(`SELECT stocks.symbol
	FROM positions INNER JOIN stocks
	ON positions.stock_id = stocks.id
	WHERE date_end IS NULL
	GROUP BY stocks.id;`)
	if err != nil {
		return
	}

	// get all stocks owned by a user
	symbols = make([]string, len(rows))
	for index, row := range rows {
		x, _ := row["symbol"]
		symbols[index] = x
	}
	return
}

func getPrices(DB *xorm.Engine, symbols []string) (err error) {
	// TODO: Is there a max batch size?? Should we batch requests in smaller portions?
	// Might want to consider using a stringbuilder for efficiency
	batchSymbols := ""
	for _, s := range symbols {
		batchSymbols += s + ","
	}

	// https://api.iextrading.com/1.0/stock/market/batch?symbols=aapl,fb&types=price
	url := "https://api.iextrading.com/1.0/stock/market/batch?symbols="
	url += batchSymbols
	url += "&types=price"

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var m map[string]interface{}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return
	}
	for key, val := range m {
		test := val.(map[string]interface{})
		for _, pps := range test {
			DB.Exec("UPDATE stocks SET price_per_share=? WHERE symbol=?", pps, key)
		}
	}
	return
}
