package model

import (
	"bytes"
	"log"
	"net/url"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/rsc/blog/atom"
)

func sanitise(feed *atom.Feed, t *atom.Text) string {
	r := strings.NewReader(t.Body)
	var w bytes.Buffer
	n, err := html.Parse(r)
	if err != nil {
		log.Print(err)
		return ""
	}
	absoluteImgTag(baseURLForFeed(feed), n)
	if err := html.Render(&w, n); err != nil {
		log.Print(err)
		return ""
	}
	return w.String()
}

// correct all <img> tags to use absolute urls.
func absoluteImgTag(base *url.URL, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "img" {
		n.Attr = stripClassAttr(n.Attr)
		for i, a := range n.Attr {
			if a.Key == "src" {
				u, err := url.Parse(a.Val)
				if err != nil {
					log.Fatal(err)
				}
				if u.Host == "" {
					// url is relative, append the path to the base url for the feed
					u = base.ResolveReference(u)
					n.Attr[i].Val = u.String()
				}
				break
			}
		}
		n.Attr = append(n.Attr, html.Attribute{Key: "class", Val: "img-responsive"})
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		absoluteImgTag(base, c)
	}
}

// remove any class attributes, we want to add out own later.
func stripClassAttr(nn []html.Attribute) []html.Attribute {
	var v []html.Attribute
	for _, n := range nn {
		if n.Key == "class" {
			continue
		}
		v = append(v, n)
	}
	return v
}

func baseURLForFeed(f *atom.Feed) *url.URL {
	for _, l := range f.Link {
		u, _ := url.Parse(l.Href)
		u.Path = "/"
		return u
	}
	panic("no base")
}
