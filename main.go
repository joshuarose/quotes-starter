package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/99designs/gqlgen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// quote represents data about random quotes
type quote struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}
type ID struct {
	ID string `json:"id"`
}

// SQL Database variable
var db *sql.DB

// Connect to Database
func connectUnixSocket() error {
	mustGetenv := func(dns string) string {
		receiveEnv := os.Getenv(dns)
		if receiveEnv == "" {
			log.Printf("Warning: %s environment variable not set.\n", dns)
		}
		return receiveEnv
	}
	// Environment variables
	var (
		dbUser         = mustGetenv("DB_USER") // postgres
		dbPwd          = mustGetenv("DB_PASS")
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET")
		dbName         = mustGetenv("DB_NAME")
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath) // SHOULD IT BE dbPASS?

	var err error
	db, err = sql.Open("pgx", dbURI) // populating package level variable with Database
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	return err
}

// Check Api Header Key
func xApiKey(c *gin.Context) bool {

	header, exists := c.Request.Header["X-Api-Key"]
	if exists {
		if header[0] == "COCKTAILSAUCE" {
			return exists
		}
	}
	return exists
}

func getRandomQuote(c *gin.Context) {
	// Check if Api Header Key exists
	if !xApiKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401 Unauthorized"})
		return
	}
	row := db.QueryRow(fmt.Sprintln("SELECT id, quote, author FROM quotes ORDER BY RANDOM() LIMIT 1"))
	q := &quote{}
	err := row.Scan(&q.ID, &q.Quote, &q.Author)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, q)

}

// Get Quote by ID
func getQuoteByIdSQL(c *gin.Context) {
	if !xApiKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401 Unauthorized"})
		return
	}
	id := c.Param("id")
	row := db.QueryRow(fmt.Sprintf("SELECT id, quote, author FROM quotes WHERE id = '%s'", id))
	q := &quote{}
	err := row.Scan(&q.ID, &q.Quote, &q.Author)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Quote ID not found"})
	} else {
		c.JSON(http.StatusOK, q)
		return
	}
}

// Create New Quote
func postNewQuote(c *gin.Context) {
	if !xApiKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401 Unauthorized"})
		return
	}

	q := &quote{} //taking JSON Body in post To insert into database
	var newID ID
	newID.ID = uuid.New().String()
	if err := c.BindJSON(&q); err != nil { //passes the HTTP status code 400 to the context and then returns a pointer or an error.
		return
	}
	// Check length of author and quote strings
	if (len(q.Quote)) < 3 || (len(q.Author)) < 3 && !xApiKey(c) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid, input must be more than three characters"})
		return
	} else {
		// Insert new Quote into Table
		sqlStatement := `INSERT INTO quotes (id, Quote, Author) VALUES ($1, $2, $3)`
		_, err := db.Exec(sqlStatement, &newID.ID, &q.Quote, &q.Author)

		if err != nil {
			fmt.Println("Error Error Error")
		}
	}
	c.JSON(http.StatusCreated, newID)

}

// Delete Quote
func deleteQuote(c *gin.Context) {
	if !xApiKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401 Unauthorized"})
		return
	}
	id := c.Param("id")
	row := db.QueryRow(fmt.Sprintf("DELETE from quotes where id = '%s'", id))
	q := &quote{}
	err := row.Scan(&q.ID, &q.Quote, &q.Author)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusNoContent, q)

}

func main() {

	err := connectUnixSocket() // call database function
	if err != nil {
		log.Println(err)
	} // Stop program if database connection fails

	rand.Seed(time.Now().UnixNano())
	router := gin.Default()
	router.GET("/quotes", getRandomQuote)
	router.GET("/quotes/:id", getQuoteByIdSQL) // ????????
	router.POST("/quotes", postNewQuote)
	router.DELETE("/quotes/:id", deleteQuote)
	router.Run("0.0.0.0:8080")

}
