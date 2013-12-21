package main

import (
	"flag"
	"html/template"
	"log"
	"path/filepath"
	"sort"
	"time"

	"github.com/pkg/math"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"

	"github.com/davecheney/planetgolang/fetch"
	"github.com/dustin/go-humanize"

	"code.google.com/p/rsc/blog/atom"
)

var (
	staticDir   = flag.String("static", filepath.Join(mustCwd(), "static"), "static asset directory")
	templateDir = flag.String("template", filepath.Join(mustCwd(), "templates"), "template directory")
)

func init() { flag.Parse() }

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
			body := sanitise(entry.Content)
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

	feeds := fetch.LoadAll(flag.Args()...)
	entries := entries(feeds)
	sort.Sort(entriesByTime(entries))
	entries = entries[:math.Max(len(entries), 10)]

	m.Get("/index", func(r render.Render) {
		s := struct {
			Title   string
			Entries []*Entry
			Feeds   []*atom.Feed
		}{"Planet Golang", entries, feeds}
		r.HTML(200, "index", &s)
	})

	m.Run()
}

type entriesByTime []*Entry

func (s entriesByTime) Len() int           { return len(s) }
func (s entriesByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s entriesByTime) Less(i, j int) bool { return s[i].Time.After(s[j].Time) }
