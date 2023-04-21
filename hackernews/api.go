package hackernews

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fastjson"
)

const (
	STORIES_URL   = "https://hacker-news.firebaseio.com/v0/topstories.json"
	ITEM_URL_BASE = "https://hacker-news.firebaseio.com/v0/item"
)

func createClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 30 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(10 * time.Second),
	}
}

func fetchJSON(url string) (*fastjson.Value, error) {
	client := createClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	value, err := fastjson.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return value, nil
}

func downloadStory(id int, wg *sync.WaitGroup, m *sync.Mutex, stories map[int]Story) {
	defer wg.Done()

	url := fmt.Sprintf("%s/%d.json", ITEM_URL_BASE, id)

	value, err := fetchJSON(url)
	if err != nil {
		log.Printf("Error fetching story with id %d: %v", id, err)
		return
	}

	story := Story{
		Title:       string(value.GetStringBytes("title")),
		By:          string(value.GetStringBytes("by")),
		Descendants: int(value.GetInt("descendants")),
		Id:          int(value.GetInt("id")),
		Score:       int(value.GetInt("score")),
		Time:        int(value.GetInt("time")),
		Type:        string(value.GetStringBytes("type")),
		URL:         string(value.GetStringBytes("url")),
	}

	kids := value.GetArray("kids")
	for i := range kids {
		story.Kids = append(story.Kids, int(kids[i].GetInt()))
	}

	m.Lock()
	stories[id] = story
	m.Unlock()
}

func FetchTopStories(number int, showSourceUrls bool) {
	value, err := fetchJSON(STORIES_URL)
	if err != nil {
		log.Printf("Error fetching top stories: %v", err)
		return
	}

	var ids []int
	valueArray := value.GetArray()
	for i := 0; i < len(valueArray); i++ {
		id := int(valueArray[i].GetInt())
		ids = append(ids, id)
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
		fmt.Print(story.PrintStyling(showSourceUrls))
	}
}
