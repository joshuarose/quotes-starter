package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"math/rand"

	"time"

	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type quote struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type newQuotes struct {
	ID string `json:"id"`
}

var quotess = map[uuid.UUID]quote{
	getUUID(): {ID: "", Quote: "Reflection is never clear.", Author: "Joe Burrow"},
	getUUID(): {ID: "", Quote: "Don't just check errors, handle them gracefully.", Author: "Oprah"},
	getUUID(): {ID: "", Quote: "A little copying is better than a little dependency.", Author: "Vienna Erhart"},
	getUUID(): {ID: "", Quote: "The bigger the interface, the weaker the abstraction.", Author: "Josh Rose"},
	getUUID(): {ID: "", Quote: "Don't panic.", Author: "Queen of England"},
}

var dbPool *sql.DB

func getUUID() uuid.UUID {
	return uuid.New()
}

var finalQuotess = map[uuid.UUID]quote{}

func setIDs() {
	for key, element := range quotess {
		element.ID = uuid.UUID.String(key)
		finalQuotess[key] = element
	}
}

var arrayOfUUIDs = []uuid.UUID{}

func makeArray() {
	for k, _ := range finalQuotess {
		arrayOfUUIDs = append(arrayOfUUIDs, k)
	}
}

func getRandomQuote() quote {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(arrayOfUUIDs))
	randomUUID := arrayOfUUIDs[randomNum]
	return finalQuotess[randomUUID]
}

func manageHeader(c *gin.Context) bool {
	headers := c.Request.Header
	header, exists := headers["X-Api-Key"]
	// fmt.Println(header)

	if exists {
		if header[0] == "COCKTAILSAUCE" {
			return true
		}
	}
	return false
}

func getQuoteByIDSQL(c *gin.Context) {
	if manageHeader(c) {
		id := c.Param("id")
		row := dbPool.QueryRow(fmt.Sprintf("select uuidkey, quote, author from quotes where uuidkey = '%s'", id))
		q := &quote{}
		err := row.Scan(&q.ID, &q.Quote, &q.Author)
		if err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusOK, q)
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func postQuoteSQL(c *gin.Context) {
	if manageHeader(c) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Print(err)
		}
		var postQuote quote
		json.Unmarshal(jsonData, &postQuote)

		sqlStatement := `
INSERT INTO quotes (uuidkey, quote, author)
VALUES ($1, $2, $3)`
		_, err = dbPool.Exec(sqlStatement, postQuote.ID, postQuote.Quote, postQuote.Author)
		formattedID := postQuote.ID
		var returnThisUUID = []newQuotes{
			{ID: formattedID},
		}
		c.JSON(http.StatusCreated, returnThisUUID[0])
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "message: An error occurred")
		}
	}
}

func main() {

	connectUnixSocket()

	setIDs()

	makeArray()

	router := gin.Default()

	router.GET("/quotes", getQuotes)
	// router.GET("/quotes/:id", getQuoteById)
	router.GET("/quotes/:id", getQuoteByIDSQL)
	// router.POST("/quotes", postQuotes)
	router.POST("/quotes", postQuoteSQL)

	// if err := dbPool.Ping(); err != nil {
	// 	log.Fatalf("unable to reach database: %v", err)
	// }
	// fmt.Println("database is reachable")

	router.Run("0.0.0.0:8080")

}

//

func getQuotes(c *gin.Context) {
	keySlice, exists := c.Request.Header["X-Api-Key"]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401"})
	} else if keySlice[0] == "COCKTAILSAUCE" {
		c.JSON(http.StatusOK, getRandomQuote())
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401"})
	}
}

// connectUnixSocket initializes a Unix socket connection pool for
// a Cloud SQL instance of Postgres.
func connectUnixSocket() error {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PWD")               // e.g. 'my-db-password'
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
	)
	fmt.Println()
	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	// fmt.Println(dbURI)

	// dbPool is the pool of database connections.
	var err error
	dbPool, err = sql.Open("pgx", dbURI)
	if err == nil {
		fmt.Println("no error")
		return fmt.Errorf("sql.Open: %v", err)

	}
	return err
}

// func getQuoteById(c *gin.Context) {
// 	keySlice, exists := c.Request.Header["X-Api-Key"]
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"status": "401"})
// 	} else if keySlice[0] == "COCKTAILSAUCE" {
// 		id := c.Param("id")
// 		for k, v := range finalQuotess {
// 			if uuid.UUID.String(k) == id {
// 				c.JSON(http.StatusOK, v)
// 				return
// 			}
// 		}
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "quote not found"})
// 	} else {
// 		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401"})
// 	}
// }

//func postQuotes(c *gin.Context) {
// 	keySlice, exists := c.Request.Header["X-Api-Key"]
// 	if !exists {
// 		c.IndentedJSON(http.StatusUnauthorized, gin.H{"status": "401"})
// 	} else if keySlice[0] == "COCKTAILSAUCE" {
// 		var newQuote quote
// 		if err := c.BindJSON(&newQuote); err != nil {
// 			return
// 		}
// 		authorCheck := len(newQuote.Author) > 3
// 		quoteCheck := len(newQuote.Quote) > 3

// 		if !authorCheck || !quoteCheck {
// 			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "400"})
// 			return
// 		}
// 		newUUID := getUUID()
// 		newID := uuid.UUID.String(newUUID)
// 		finalQuotess[newUUID] = newQuote
// 		newQuote.ID = newID
// 		var UseThisUUID = []newQuotes{
// 			{ID: newQuote.ID},
// 		}
// 		c.IndentedJSON(http.StatusCreated, UseThisUUID[0])
// 		makeArray()
// 	} else {
// 		c.IndentedJSON(http.StatusUnauthorized, gin.H{"status": "401"})
// 	}
// }
