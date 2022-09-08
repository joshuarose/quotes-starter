package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// goQuotes represents data about random quotes
type goQuotes struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

var randomQuotes = []goQuotes{
	{Quote: "Don't communicate by sharing memory, share memory by communicating.", Author: "Rob Pike"},
	{Quote: "With the unsafe package there are no guarantees.", Author: "Rob Pike"},
	{Quote: "A little copying is better than a little dependency.", Author: "Rob Pike"},
	{Quote: "Design the architecture, name the components, document the details.", Author: "Rob Pike"},
	{Quote: "Don't just check errors, handle them gracefully.", Author: "Rob Pike"},
	{Quote: "Avoid unused method receiver names", Author: "Kalese Carpenter"},
}

func getRandomQuotes(context *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(len(randomQuotes))
	context.JSON(http.StatusOK, randomQuotes[randomNumber])
}

// goQuotes := (randomQuotes[randomNumber])

func main() {
	router := gin.Default()
	router.GET("/quotes", getRandomQuotes)
	router.Run("0.0.0.0:8080")

}
