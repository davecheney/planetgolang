package model

import (
	"html/template"
	"log"
	"sort"
	"sync"
	"time"

	"code.google.com/p/rsc/blog/atom"
)

type Model struct {
	mu      sync.Mutex
	entries []*Entry
	feeds	[]*atom.Feed
}

func New(urls []string) *Model {
	var m Model
	go func() {
		feeds := m.LoadAll(urls)
		entries := entries(feeds)
		sort.Sort(entriesByTime(entries))
		m.mu.Lock()
		m.entries = entries
		m.feeds = feeds
		m.mu.Unlock()
	}()
	return &m
}

func (m *Model) Entries() []*Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.entries
}

func (m *Model) Feeds() []*atom.Feed {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.feeds
}

func (m *Model) LoadAll(urls []string) []*atom.Feed {
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

type Entry struct {
	*atom.Feed
	*atom.Entry
	time.Time
	Content template.HTML
}

func entries(feeds []*atom.Feed) []*Entry {
	var entries []*Entry
	for _, feed := range feeds {
		for _, entry := range feed.Entry {
			t := saneDate(entry)
			body := sanitise(feed, entry.Content)
			entries = append(entries, &Entry{feed, entry, t, template.HTML(body)})
		}
	}
	return entries
}

func saneDate(entry *atom.Entry) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05-07:00", string(entry.Published))
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", string(entry.Published))
		if err != nil {
			log.Print(err)
		}
	}
	return t
}

type entriesByTime []*Entry

func (s entriesByTime) Len() int           { return len(s) }
func (s entriesByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s entriesByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
