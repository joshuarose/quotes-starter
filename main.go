package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"math/rand"
)

type quote struct {
	Quote string `json:"quote"`
	Author string `json:"author"`
}

var quotes = []quote{
	{Quote: "Clear is better than clever.", Author: "Ronald McDonald"},
	{Quote: "Empty string check!", Author: "Squidward Tentacles"},
	{Quote: "Don't panic.", Author: "Oprah Winfrey"},
	{Quote: "A little copying is better than a little dependency.", Author: "Chris Pratt"},
	{Quote: "The bigger the interface, the weaker the abstraction.", Author: "Mary Poppins"},
	{Quote: "With the unsafe package there are no guarantees.", Author: "Rob Dyrdek"},
	{Quote: "Reflection is never clear.", Author: "Bobby Hill"},
	{Quote: "Don't just check errors, handle them gracefully.", Author: "Shrek"},
	{Quote: "Documentation is for users.", Author: "Hermione Granger"},
	{Quote: "Errors are values.", Author: "Clark Kent"},
	{Quote: "Make the zero value useful.", Author: "Drake"},
	{Quote: "Channels orchestrate; mutexes serialize.", Author: "Yo-Yo Ma"},
	{Quote: "Don't communicate by sharing memory, share memory by communicating.", Author: "Prince"},
	{Quote: "Concurrency is not parallelism.", Author: "Lao Tzu"},
	{Quote: "Design the architecture, name the components, document the details.", Author: "Tony the Tiger"},
}

func main() {
	router := gin.Default()
	router.GET("/quotes", getRandomQuote)
	router.Run("0.0.0.0:8080")
}

func getRandomQuote(c *gin.Context) {
	index := rand.Intn(len(quotes))
	quote := quotes[index]
	c.IndentedJSON(http.StatusOK, quote)
}

// func getQuotes(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, quotes)
// }

// func main() {
// 	r := gin.Default()
// 	r.GET("/quotes", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"quote": getQuote(),
// 			"author": "Jessie & Nate - Gold Team Rules!",
// 		})
// 	})
// 	r.Run() //default's on 8080
// }