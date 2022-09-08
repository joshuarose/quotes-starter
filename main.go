package main

import (
	"fmt"
)

// goQuotes represents data about random quotes
type goQuotes struct {
	ID     string `json:"id"`
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

var randomQuotes = []goQuotes{
	{ID: "1", Quote: "Don't communicate by sharing memory, share memory by communicating.", Author: "Rob Pike"},
	{ID: "2", Quote: "With the unsafe package there are no guarantees.", Author: "Rob Pike"},
	{ID: "3", Quote: "A little copying is better than a little dependency.", Author: "Rob Pike"},
	{ID: "4", Quote: "Design the architecture, name the components, document the details.", Author: "Rob Pike"},
	{ID: "5", Quote: "Don't just check errors, handle them gracefully.", Author: "Rob Pike"},
	{ID: "6", Quote: "Avoid unused method receiver names", Author: "Kalese Carpenter"},
}

func main() {
	fmt.Println(randomQuotes)
}
