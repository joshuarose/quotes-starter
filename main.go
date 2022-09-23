package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"math/rand"

	"time"

	"github.com/google/uuid"
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

	setIDs()

	makeArray()

	router := gin.Default()

	router.GET("/quotes", getQuotes)
	router.GET("/quotes/:id", getQuoteById)
	router.POST("/quotes", postQuotes)

	router.Run("0.0.0.0:8080")
}

func postQuotes(c *gin.Context) {
	keySlice := c.Request.Header["X-Api-Key"]
	keyString := keySlice[0]
	if keyString == "COCKTAILSAUCE" {
		var newQuote quote
		// Call BindJSON to bind the received JSON to
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
		finalQuotess[newUUID] = newQuote
		newQuote.ID = uuid.UUID.String(newUUID)
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
	keySlice := c.Request.Header["X-Api-Key"]
	keyString := keySlice[0]
	if keyString == "COCKTAILSAUCE" {
		c.JSON(http.StatusOK, getRandomQuote())
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "401"})
	}
}

func getQuoteById(c *gin.Context) {
	keySlice := c.Request.Header["X-Api-Key"]
	keyString := keySlice[0]
	if keyString == "COCKTAILSAUCE" {
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
