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

var dbPool *sql.DB

func getUUID() uuid.UUID {
	return uuid.New()
}

func manageHeader(c *gin.Context) bool {
	headers := c.Request.Header
	header, exists := headers["X-Api-Key"]
	if exists {
		if header[0] == "COCKTAILSAUCE" {
			return true
		}
	}
	return false
}

func postQuote(c *gin.Context) {
	if manageHeader(c) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, "message: error")
		}
		var postQuote quote
		json.Unmarshal(jsonData, &postQuote)
		if len(postQuote.Quote) < 3 || len(postQuote.Author) < 3 {
			c.JSON(http.StatusBadRequest, "message: Quote and Author must exceed 3 characters")
			return
		}
		sqlStatement := `
INSERT INTO quotes (uuidkey, quote, author)
VALUES ($1, $2, $3)`
		generatedUUID := getUUID()
		formattedID := generatedUUID.String()
		_, err = dbPool.Exec(sqlStatement, formattedID, postQuote.Quote, postQuote.Author)
		var returnThisUUID = []newQuotes{
			{ID: formattedID},
		}
		c.JSON(http.StatusCreated, returnThisUUID[0])
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "message: An error occurred")
			return
		}
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func deleteQuotesByID(c *gin.Context) {
	if manageHeader(c) {
		id := c.Param("id")
		statement := `DELETE FROM quotes WHERE uuidkey=$1 `
		_, err2 := dbPool.Exec(statement, id)
		if err2 != nil {
			log.Println(err2)
			return
		}
		c.JSON(http.StatusNoContent, "message: Successfully deleted")
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func getRandomQuote(c *gin.Context) {
	if manageHeader(c) {
		row := dbPool.QueryRow(fmt.Sprintln("select uuidkey, quote, author from quotes order by random() limit 1"))
		q := &quote{}
		err := row.Scan(&q.ID, &q.Quote, &q.Author)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, q)
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func getQuoteByID(c *gin.Context) {
	if manageHeader(c) {
		id := c.Param("id")
		row := dbPool.QueryRow(fmt.Sprintf("select uuidkey, quote, author from quotes where uuidkey = '%s'", id))
		q := &quote{}
		switch err := row.Scan(&q.ID, &q.Quote, &q.Author); err {
		case sql.ErrNoRows:
			c.JSON(http.StatusNotFound, "message: ID does not exist")
			return
		case nil:
			c.JSON(http.StatusOK, q)
			return
		default:
			c.JSON(http.StatusNotFound, "message: Something went wrong")
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}

func main() {
	defer dbPool.Close()
	connectUnixSocket()

	router := gin.Default()

	router.GET("/quotes/:id", getQuoteByID)
	router.POST("/quotes", postQuote)
	router.DELETE("/quotes/:id", deleteQuotesByID)
	router.GET("/quotes", getRandomQuote)
	router.Run("0.0.0.0:8080")
}

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

	var err error
	dbPool, err = sql.Open("pgx", dbURI)
	if err == nil {
		fmt.Println("no error")
		return fmt.Errorf("sql.Open: %v", err)

	}
	return err
}
