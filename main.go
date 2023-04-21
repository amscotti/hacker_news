package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/amscotti/hacker_news/hackernews"
)

func main() {
	var number int
	var showSourceUrls bool

	flag.IntVar(&number, "n", 5, "Specify the number of top stories to display (default: 5).")
	flag.BoolVar(&showSourceUrls, "u", false, "Include the source URLs of the stories in the output.")
	flag.Parse()

	if number <= 0 {
		err := errors.New("number of stories must be a positive integer")
		fmt.Println(err)
		return
	}

	hackernews.FetchTopStories(number, showSourceUrls)
}
