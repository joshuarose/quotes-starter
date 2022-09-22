package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// goQuotes represents data about random quotes
type goQuote struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

var randomQuote = map[string]goQuote{
	"374be3f1-956a-4169-874a-0632c09a2599": {ID: "374be3f1-956a-4169-874a-0632c09a2599", Quote: "Don't communicate by sharing memory, share memory by communicating.", Author: "Rob Pike"},
	"a4539044-da8d-4064-bb05-2421abd4c77d": {ID: "a4539044-da8d-4064-bb05-2421abd4c77d", Quote: "With the unsafe package there are no guarantees.", Author: "Rob Pike"},
	"068faa87-9afa-4f7f-8aed-ff2d303c79e5": {ID: "068faa87-9afa-4f7f-8aed-ff2d303c79e5", Quote: "A little copying is better than a little dependency.", Author: "Rob Pike"},
	"0f4036b0-d49a-46b9-9ec2-577fbfd4f714": {ID: "0f4036b0-d49a-46b9-9ec2-577fbfd4f714", Quote: "Design the architecture, name the components, document the details.", Author: "Rob Pike"},
	"10a2781c-113f-4c49-a670-8ed322882f1a": {ID: "10a2781c-113f-4c49-a670-8ed322882f1a", Quote: "Don't just check errors, handle them gracefully.", Author: "Rob Pike"},
	"77efbc8b-2289-45ee-9461-b1f602fecf3e": {ID: "77efbc8b-2289-45ee-9461-b1f602fecf3e", Quote: "Avoid unused method receiver names", Author: "Kalese Carpenter"},
	"211cf4f3-3893-43b8-a1d2-88aedc14df5a": {ID: "211cf4f3-3893-43b8-a1d2-88aedc14df5a", Quote: "Gofmt's style is no one's favorite, yet gofmt is everyone's favorite", Author: "Rob Pike"},
	"323d8e20-7975-4ff1-af6d-99dc7f57f35a": {ID: "323d8e20-7975-4ff1-af6d-99dc7f57f35a", Quote: "For brands or words with more than 1 capital letter, lowercase all letters", Author: "Kalese Carpenter"},
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := gin.Default()
	router.GET("/quotes", getRandomQuotes)
	router.GET("/quotes/:id", getQuoteById)
	router.Run("0.0.0.0:8080")
}

// Get A Random Quote From Map
func getRandomQuotes(c *gin.Context) {
	counter := 0
	randomNumber := rand.Intn(len(randomQuote))

	for _, v := range randomQuote {
		if counter == randomNumber {
			c.JSON(http.StatusOK, &v)
		}
		counter++
	}
}

// Get Quote By ID
func getQuoteById(c *gin.Context) {
	id := c.Param("id")
	singleQuote, exists := randomQuote[id]
	if exists {
		c.JSON(http.StatusOK, singleQuote)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "404 Not Found"}) // change this to 404 message
}

//Post New Quote
