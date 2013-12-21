package main

import (
	"bytes"
	"log"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/rsc/blog/atom"
)

func sanitise(t *atom.Text) string {
	log.Printf("Processing content type: %q", t.Type)
	r := strings.NewReader(t.Body)
	var w bytes.Buffer
	n, err := html.Parse(r)
	if err != nil {
		log.Print(err)
		return ""
	}
	if err := html.Render(&w, n); err != nil {
		log.Print(err)
		return ""
	}
	return w.String()
}
