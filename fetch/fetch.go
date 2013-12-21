package fetch

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"sync"

	"code.google.com/p/rsc/blog/atom"
)

func load(url string) (*atom.Feed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("non 200 response code")
	}
	var feed atom.Feed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, errors.New("non 200 response code")
	}
	return &feed, nil
}

func LoadAll(urls ...string) []*atom.Feed {
	var wg sync.WaitGroup
	var feeds []*atom.Feed
	var c = make(chan *atom.Feed, len(urls))
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			f, err := load(url)
			if err != nil {
				log.Print(err)
			} else {
				c <- f
			}
		}(url)
	}
	wg.Wait()
	close(c)
	for f := range c {
		feeds = append(feeds, f)
	}
	return feeds
}
