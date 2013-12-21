package main

import (
	"flag"
	"html/template"
	"path/filepath"

	"github.com/pkg/math"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/davecheney/planetgolang/model"
	"github.com/dustin/go-humanize"

	"code.google.com/p/rsc/blog/atom"
)

const ENTRIES_PER_PAGE = 25

var (
	staticDir   = flag.String("static", filepath.Join(mustCwd(), "static"), "static asset directory")
	templateDir = flag.String("template", filepath.Join(mustCwd(), "templates"), "template directory")
)

func init() { flag.Parse() }

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
			"humanize": humanize.Time,
			"url": func(l []atom.Link) string {
				for _, l := range l {
					return l.Href
				}
				return "#"
			},
		}},
		Charset: "utf-8",
	}))

	mod := model.New(flag.Args())

	m.Get("/index", func(r render.Render) {
		entries := mod.Entries()
		entries = entries[:math.Min(len(entries), ENTRIES_PER_PAGE)]
		s := struct {
			Title   string
			Entries []*model.Entry
			Feeds   []*atom.Feed
		}{"Planet Golang", entries, mod.Feeds()}
		r.HTML(200, "index", &s)
	})

	m.Run()
}
