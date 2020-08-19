package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	_ "github.com/apache/calcite-avatica-go"
)

/* TODO:
- Bug fixes/hardening
- Routing back from single stock page
- UI refresh?

5. Bonus:
Fetch Tweets
Store tweets in OpDB
Render tweets
*/

/*
{
	"Global Quote": {
		"01. symbol": "CLDR",
		"02. open": "11.3200",
		"03. high": "11.5000",
		"04. low": "11.2933",
		"05. price": "11.3800",
		"06. volume": "1647629",
		"07. latest trading day": "2020-08-17",
		"08. previous close": "11.2600",
		"09. change": "0.1200",
		"10. change percent": "1.0657%"
	}
}
*/
type globalQuote struct {
	Quote quote `json:"Global Quote"`
}

type quote struct {
	Symbol  string  `json:"01. symbol"`
	Price   float32 `json:"05. price,string"`
	ChangeP string  `json:"10. change percent"`
}

type cleanquote struct {
	Symbol  string  `json:"symbol"`
	Price   float32 `json:"price"`
	ChangeP string  `json:"changep"`
}

var db *sql.DB

func main() {
	url := "http://localhost:8765/"

	var err error
	db, err = sql.Open("avatica", url)
	if err != nil {
		log.Fatal("Connection: ", err)
	}
	defer db.Close()

	//dropAndCreateDBTable(db)

	parseTS("IBM")
	parseTS("CLDR")
	parseTS("MSFT")

	putTickerDB(fetchTicker("msft"), db)
	putTickerDB(fetchTicker("IBM"), db)
	putTickerDB(fetchTicker("CLDR"), db)

	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.GET("/ticker/:tickerID", TickerHandler)
	api.GET("/list", ListHandler)
	api.GET("/intraday/:tickerID", IntraHandler)

	// Start and run the server
	router.Run(":3000")

}

func fetchTicker(ticker string) (q quote) {
	apikey := os.Getenv("AV_API_KEY")
	url := "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + ticker + "&apikey=" + apikey
	avClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "demo")

	res, getErr := avClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	GQuote := globalQuote{}
	jsonErr := json.Unmarshal(body, &GQuote)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return GQuote.Quote
}

func dropAndCreateDBTable(db *sql.DB) {
	log.Println("Deleting existing table...")
	_, err := db.Exec("DROP TABLE IF EXISTS stocks")
	if err != nil {
		log.Fatal("Could not drop table", err)
	}

	log.Println("Create table if not exists...")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stocks (symbol VARCHAR PRIMARY KEY, price FLOAT, changep VARCHAR)")
	if err != nil {
		log.Fatal("Create: ", err)
	}
}

func putTickerDB(q quote, db *sql.DB) {
	log.Println("Upsert rows...")
	_, err := db.Exec("UPSERT INTO stocks VALUES (?, ?, ?)", q.Symbol, q.Price, q.ChangeP)
	if err != nil {
		log.Println(err)
	}

}

func getTickerDB(symbol string, db *sql.DB) (q quote) {
	qrow := quote{}
	rows, err := db.Query("SELECT symbol, price, changep FROM stocks WHERE symbol = '" + symbol + "'")
	if err != nil {
		log.Fatal("Query: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&qrow.Symbol, &qrow.Price, &qrow.ChangeP)
		if err != nil {
			log.Fatal(err)
		}
	}
	return qrow
}

func getTickersDB(count int, db *sql.DB) (qs []cleanquote) {
	var res []cleanquote
	qrow := cleanquote{}
	rows, err := db.Query("SELECT symbol, price, changep FROM stocks LIMIT 10")
	if err != nil {
		log.Fatal("Query: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&qrow.Symbol, &qrow.Price, &qrow.ChangeP)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, qrow)
	}
	return res
}

func ListHandler(c *gin.Context) {
	//res := []cleanquote{{"CLDR", 12.40}}
	res := getTickersDB(10, db)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, res)
}

func TickerHandler(c *gin.Context) {
	q := getTickerDB(strings.ToUpper(c.Param("tickerID")), db)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"symbol":  q.Symbol,
		"price":   q.Price,
		"changep": q.ChangeP,
	})
}

func fetchTS(ticker string) (r []byte) {
	apikey := os.Getenv("AV_API_KEY")
	url := "https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=" + ticker + "&interval=5min&apikey=" + apikey
	avClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "demo")

	res, getErr := avClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//	GQuote := globalQuote{}
	//	jsonErr := json.Unmarshal(body, &GQuote)
	//	if jsonErr != nil {
	//		log.Fatal(jsonErr)
	//	}

	return body
}

func parseTS(symbol string) {
	log.Println("Creating values table if not exists...")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS stockvals (symbol VARCHAR NOT NULL, instant TIMESTAMP NOT NULL, price FLOAT, CONSTRAINT pk PRIMARY KEY(symbol,instant))")
	if err != nil {
		log.Fatal("Create: ", err)
	}

	testJson := fetchTS(symbol)
	var result map[string]interface{}
	json.Unmarshal([]byte(testJson), &result)
	times := result["Time Series (5min)"].(map[string]interface{})

	log.Println("Upsert value rows...")
	for key, value := range times {
		vals := value.(map[string]interface{})

		closeprice, err := strconv.ParseFloat(vals["4. close"].(string), 32)
		if err != nil {
			log.Println(err)
		}
		_, err = db.Exec("UPSERT INTO stockvals VALUES (?, ?, ?)", symbol, key, closeprice)
		if err != nil {
			log.Println(err)
		}
	}
	getTSDB(symbol)
}

func getTSDB(symbol string) (times []int64, prices []float32) {
	var (
		instant string
		price   float32
	)

	rows, err := db.Query("SELECT instant, price FROM stockvals WHERE symbol = '" + symbol + "' ORDER BY instant ASC")
	if err != nil {
		log.Fatal("Query: ", err)
	}
	defer rows.Close()
	times = make([]int64, 0)
	prices = make([]float32, 0)
	for rows.Next() {
		err := rows.Scan(&instant, &price)
		if err != nil {
			log.Fatal(err)
		}
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, instant)
		if err != nil {
			fmt.Println(err)
		}
		times = append(times, t.Unix())
		prices = append(prices, price)
	}
	return times, prices
}

func IntraHandler(c *gin.Context) {
	times, prices := getTSDB(strings.ToUpper(c.Param("tickerID")))
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"times":  times,
		"prices": prices,
	})
}
