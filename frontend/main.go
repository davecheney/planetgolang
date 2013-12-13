package main

import (
	"flag"
	"path/filepath"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
)

var (
	staticDir = flag.String("static", filepath.Join(mustCwd(), "static"), "static asset directory")
	templateDir = flag.String("template", filepath.Join(mustCwd(), "templates"), "template directory")
)

func init() { flag.Parse() }

func main() {
	m := martini.Classic()
	
	// setup static assets
	m.Use(martini.Static(*staticDir))

	// setup templates
	m.Use(render.Renderer(render.Options{
  		Directory: *templateDir, 
  		Extensions: []string{".tmpl"}, 
	})

	m.Run()
}
