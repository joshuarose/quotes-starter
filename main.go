package main

import (
	"database/sql"
	"fmt"
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

func main() {
	connectUnixSocket()

	row := dbPool.QueryRow("SELECT * FROM quotes LIMIT 1;")
	fmt.Println(row)

	setIDs()

	makeArray()

	router := gin.Default()

	router.GET("/quotes", getQuotes)
	router.GET("/quotes/:id", getQuoteById)
	router.POST("/quotes", postQuotes)

	if err := dbPool.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")

	router.Run("0.0.0.0:8080")

}

func postQuotes(c *gin.Context) {
	keySlice, exists := c.Request.Header["X-Api-Key"]
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"status": "401"})
	} else if keySlice[0] == "COCKTAILSAUCE" {
		var newQuote quote
		if err := c.BindJSON(&newQuote); err != nil {
			return
		}
		authorCheck := len(newQuote.Author) > 3
		quoteCheck := len(newQuote.Quote) > 3

		if !authorCheck || !quoteCheck {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "400"})
			return
		}
		newUUID := getUUID()
		newID := uuid.UUID.String(newUUID)
		finalQuotess[newUUID] = newQuote
		newQuote.ID = newID
		var UseThisUUID = []newQuotes{
			{ID: newQuote.ID},
		}
		c.IndentedJSON(http.StatusCreated, UseThisUUID[0])
		makeArray()
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"status": "401"})
	}
}

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

func getQuoteById(c *gin.Context) {
	keySlice, exists := c.Request.Header["X-Api-Key"]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401"})
	} else if keySlice[0] == "COCKTAILSAUCE" {
		id := c.Param("id")
		for k, v := range finalQuotess {
			if uuid.UUID.String(k) == id {
				c.JSON(http.StatusOK, v)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "quote not found"})
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "401"})
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
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PWD")               // e.g. 'my-db-password'
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
	)
	fmt.Println()
	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	fmt.Println(dbURI)

	// dbPool is the pool of database connections.
	var err error
	dbPool, err = sql.Open("pgx", dbURI)
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	// ...

	return err
}
