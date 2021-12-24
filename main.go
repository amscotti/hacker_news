package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	STORIES_URL   = "https://hacker-news.firebaseio.com/v0/topstories.json"
	ITEM_URL_BASE = "https://hacker-news.firebaseio.com/v0/item"
)

type Story struct {
	Title       string
	By          string
	Descendants int
	Id          int
	Kids        []int
	Score       int
	Time        int
	Type        string
	URL         string
}

func (s *Story) Print(showSourceUrl bool) {
	fmt.Println(s.Title)
	fmt.Printf("score: %d\tcomments: %d\tuser: %s\n", s.Score, s.Descendants, s.By)
	fmt.Printf("url: https://news.ycombinator.com/item?id=%d\n", s.Id)
	if showSourceUrl {
		fmt.Println(s.URL)
	}
	fmt.Println("")
}

func downloadStory(id int, wg *sync.WaitGroup, m *sync.Mutex, stories map[int]Story) {
	url := fmt.Sprintf("%s/%d.json", ITEM_URL_BASE, id)

	rsp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var story Story
	if err := json.Unmarshal(data, &story); err != nil {
		log.Fatal(err)
	}

	m.Lock()
	stories[id] = story
	m.Unlock()

	wg.Done()
}

func fetch(number int, showSourceUrls bool) {
	fmt.Println("Fetching latest stories...")

	rsp, err := http.Get(STORIES_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ids []int
	if err := json.Unmarshal(data, &ids); err != nil {
		log.Fatal(err)
	}

	if len(ids) > number {
		ids = ids[:number]
	}

	storyDetails := make(map[int]Story)

	var m sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(ids))

	for _, id := range ids {
		go downloadStory(id, &wg, &m, storyDetails)
	}

	wg.Wait()

	for _, id := range ids {
		story := storyDetails[id]
		story.Print(showSourceUrls)
	}
}

func main() {
	var number int
	var showSourceUrls bool

	flag.IntVar(&number, "n", 5, "Number of top news to show.")
	flag.BoolVar(&showSourceUrls, "u", false, "Show source urls.")
	flag.Parse()

	fetch(number, showSourceUrls)
}
