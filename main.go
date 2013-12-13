package main

import (
	"encoding/xml"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/dustin/go-humanize"

	"code.google.com/p/rsc/blog/atom"
)

var (
	staticDir   = flag.String("static", filepath.Join(mustCwd(), "static"), "static asset directory")
	templateDir = flag.String("template", filepath.Join(mustCwd(), "templates"), "template directory")
)

func init() { flag.Parse() }

func load() []*atom.Feed {
	var feeds []*atom.Feed
	for _, url := range flag.Args() {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != 200 {
			log.Fatal("non 200 response code")
		}
		var feed atom.Feed
		if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		feeds = append(feeds, &feed)
	}
	return feeds
}

type Entry struct {
	*atom.Feed
	*atom.Entry
	time.Time
}

func entries(feeds []*atom.Feed) []*Entry {
	var entries []*Entry
	for _, feed := range feeds {
		for _, entry := range feed.Entry {
			t, err := time.Parse("2006-01-02T15:04:05-07:00", string(entry.Published))
			if err != nil {
				//log.Fatal(err)
			}
			entries = append(entries, &Entry{feed, entry, t})
		}
	}
	return entries
}

func main() {
	m := martini.Classic()

	// setup static assets
	m.Use(martini.Static(*staticDir))

	// setup templates
	m.Use(render.Renderer(render.Options{
		Directory:  *templateDir,
		Extensions: []string{".tmpl"},
		Layout:     "layout",
		Funcs: []template.FuncMap{{
			"pp":       func(s string) template.HTML { return template.HTML(s) },
			"humanize": humanize.Time,
			"url": func(l []atom.Link) string {
				for _, l := range l {
					return l.Href
				}
				return "#"
			},
		}},
	}))

	feeds := load()
	entries := entries(feeds)
	sort.Sort(entriesByTime(entries))

	m.Get("/index", func(r render.Render) {
		s := struct {
			Title   string
			Entries []*Entry
			Feeds   []*atom.Feed
		}{"Index", entries, feeds}
		r.HTML(200, "index", &s)
	})

	m.Run()
}

type entriesByTime []*Entry

func (s entriesByTime) Len() int           { return len(s) }
func (s entriesByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s entriesByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
