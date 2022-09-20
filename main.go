package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"

	"math/rand"

	"time"
)

type quote struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

var quotes = []quote{
	{ID: "1", Quote: "Reflection is never clear.", Author: "Joe Burrow"},
	{ID: "2", Quote: "Don't just check errors, handle them gracefully.", Author: "Oprah"},
	{ID: "3", Quote: "A little copying is better than a little dependency.", Author: "Vienna Erhart"},
	{ID: "4", Quote: "The bigger the interface, the weaker the abstraction.", Author: "Josh Rose"},
	{ID: "5", Quote: "Don't panic.", Author: "Queen of England"},
}

func getRandomQuote() quote {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(quotes))
	randomQuote := quotes[randomNum]
	return randomQuote
}

func main() {
	fmt.Print(getRandomQuote())

	router := gin.Default()

	router.GET("/quotes", getQuotes)
	router.GET("/quotes/:id", getQuoteById)
	router.POST("/quotes", postQuotes)
	// router.DELETE("/quotes/:id", deleteQuotes)

	router.Run("0.0.0.0:8080")
}

// postAlbums adds an album from JSON received in the request body.
func postQuotes(c *gin.Context) {
	var newQuote quote

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newQuote); err != nil {
		return
	}

	// Add the new album to the slice.
	quotes = append(quotes, newQuote)
	c.IndentedJSON(http.StatusCreated, newQuote)
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
	id := c.Param("id")
	for _, a := range quotes {
		if a.ID == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "quote not found"})
}

// func deleteQuotes(c *gin.Context) {
// 	c.JSON(http.StatusOK, getRandomQuote())
// for (i=0; i <= len(quotes); i++) {

// }
// }
