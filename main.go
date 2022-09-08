package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type quote struct {
	ID     string `json:"id"`
	Quote  string `json:"words"`
	Author string `json:"author"`
}

var quotes = []quote{
	{ID: "1", Quote: "Reflection is never clear.", Author: "Joe Burrow"},
	{ID: "2", Quote: "Don't just check errors, handle them gracefully.", Author: "Oprah"},
	{ID: "3", Quote: "A little copying is better than a little dependency.", Author: "Vienna Erhart"},
	{ID: "4", Quote: "The bigger the interface, the weaker the abstraction.", Author: "Josh Rose"},
	{ID: "5", Quote: "Don't panic.", Author: "Queen of England"},
}

//randomQuote :=

func main() {
	router := gin.Default()
	router.GET("/quotes", getQuotes)
	router.Run("localhost:8080")
}

func getQuotes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, quotes)
}
