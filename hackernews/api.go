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

type StoryIdPair struct {
	PlacementId int
	StoryId     int
}

var clientPool = sync.Pool{
	New: func() interface{} {
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     30 * time.Second,
		}
		return &http.Client{
			Transport: transport,
			Timeout:   time.Duration(10 * time.Second),
		}
	},
}

func fetchJSON(url string) (*fastjson.Value, error) {
	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

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

func downloadStory(storyIdChan <-chan StoryIdPair, storyChan chan<- Story) {
	for storyId := range storyIdChan {
		url := fmt.Sprintf("%s/%d.json", ITEM_URL_BASE, storyId.StoryId)

		value, err := fetchJSON(url)
		if err != nil {
			log.Printf("Error fetching story with id %d: %v", storyId.StoryId, err)
			storyChan <- Story{}
			return
		}

		story := Story{
			Index:       storyId.PlacementId,
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

		storyChan <- story
	}
}

func FetchTopStories(number int, showSourceUrls bool) {
	value, err := fetchJSON(STORIES_URL)
	if err != nil {
		log.Printf("Error fetching top stories: %v", err)
		return
	}

	valueArray := value.GetArray()
	count := min(len(valueArray), number)

	stories := make([]Story, count)

	storyIdChan := make(chan StoryIdPair, count)
	storyChan := make(chan Story, count)

	for range count {
		go downloadStory(storyIdChan, storyChan)
	}

	go func() {
		for index, id := range valueArray[:number] {
			storyIdChan <- StoryIdPair{PlacementId: index, StoryId: id.GetInt()}
		}
		close(storyIdChan)
	}()

	for range count {
		story := <-storyChan
		stories[story.Index] = story
	}

	close(storyChan)

	for _, story := range stories {
		fmt.Print(story.PrintStyling(showSourceUrls))
	}
}
