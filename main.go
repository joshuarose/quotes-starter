package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

//  Quotes Map

var quotesMap = map[string]quote{
	"374be3f1-956a-4169-874a-0632c09a2599": {ID: "374be3f1-956a-4169-874a-0632c09a2599", Quote: "Don't communicate by sharing memory, share memory by communicating.", Author: "Rob Pike"},
	"a4539044-da8d-4064-bb05-2421abd4c77d": {ID: "a4539044-da8d-4064-bb05-2421abd4c77d", Quote: "With the unsafe package there are no guarantees.", Author: "Rob Pike"},
	"068faa87-9afa-4f7f-8aed-ff2d303c79e5": {ID: "068faa87-9afa-4f7f-8aed-ff2d303c79e5", Quote: "A little copying is better than a little dependency.", Author: "Rob Pike"},
	"0f4036b0-d49a-46b9-9ec2-577fbfd4f714": {ID: "0f4036b0-d49a-46b9-9ec2-577fbfd4f714", Quote: "Design the architecture, name the components, document the details.", Author: "Rob Pike"},
	"10a2781c-113f-4c49-a670-8ed322882f1a": {ID: "10a2781c-113f-4c49-a670-8ed322882f1a", Quote: "Don't just check errors, handle them gracefully.", Author: "Rob Pike"},
	"77efbc8b-2289-45ee-9461-b1f602fecf3e": {ID: "77efbc8b-2289-45ee-9461-b1f602fecf3e", Quote: "Avoid unused method receiver names", Author: "Kalese Carpenter"},
	"211cf4f3-3893-43b8-a1d2-88aedc14df5a": {ID: "211cf4f3-3893-43b8-a1d2-88aedc14df5a", Quote: "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite", Author: "Rob Pike"},
	"323d8e20-7975-4ff1-af6d-99dc7f57f35a": {ID: "323d8e20-7975-4ff1-af6d-99dc7f57f35a", Quote: "For brands or words with more than 1 capital letter, lowercase all letters", Author: "Kalese Carpenter"},
}

var db *sql.DB

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
	router.DELETE("/quotes", deleteQuote)
	router.Run("0.0.0.0:8080")

}

// Connect to Database
func connectUnixSocket() error {
	mustGetenv := func(dns string) string {
		receiveEnv := os.Getenv(dns)
		if receiveEnv == "" {
			log.Printf("Warning: %s environment variable not set.\n", dns)
		}
		return receiveEnv
	}

	var (
		dbUser         = mustGetenv("DB_USER") // postgres
		dbPwd          = mustGetenv("DB_PASS")
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET")
		dbName         = mustGetenv("DB_NAME")
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath) // SHOULD IT BE dbPASS?

	var err error
	db, err = sql.Open("pgx", dbURI) // populating package level variable with DB
	if err != nil {
		return fmt.Errorf("sql.Open: %v", err)
	}

	return err
}

// Get A Random Quote From Map
func getRandomQuote(c *gin.Context) {
	// Check if Api Header Key exists
	if !xApiKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401 Unauthorized"})
		return
	}
	counter := 0
	randomNumber := rand.Intn(len(quotesMap))

	for _, v := range quotesMap {
		if counter == randomNumber {
			c.JSON(http.StatusOK, &v)
		}
		counter++
	}
}

// Get Quote by ID
func getQuoteByIdSQL(c *gin.Context) {
	if xApiKey(c) {
		id := c.Param("id")
		row := db.QueryRow(fmt.Sprintf("select id, quote, author from quotes where id = '%s'", id))
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

// Create New Quote
func postNewQuote(c *gin.Context) {
	// Check if Api Header Key exists
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
	// Insert new Quote into Table
	sqlStatement := `INSERT INTO quotes (id, Quote, Author) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, &newID.ID, &q.Quote, &q.Author)

	if err != nil {
		fmt.Println("Error Error Error")
	}

	c.JSON(http.StatusCreated, newID)

	if (len(q.Quote)) < 3 || (len(q.Author)) < 3 { // Check length of author and quote strings
		c.JSON(http.StatusBadRequest, gin.H{"status": "400 Bad Request"})
		return
	}

	c.JSON(http.StatusCreated, newID)
}

// Get Api Header Key
func xApiKey(c *gin.Context) bool {

	header, exists := c.Request.Header["X-Api-Key"]
	if exists {
		if header[0] == "COCKTAILSAUCE" {
			return exists
		}
	}
	return exists
}

// Delete Quote
func deleteQuote(c *gin.Context) {
	if xApiKey(c) {
		id := c.Param("id")
		row := db.QueryRow(fmt.Sprintf("delete from quotes where id = '%s'", id))
		q := &quote{}
		err := row.Scan(&q.ID, &q.Quote, &q.Author)
		if err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusNoContent, q)
		return
	}

}
