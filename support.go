package main

import (
	"log"
	"os"

	"code.google.com/p/rsc/blog/atom"
)

func mustCwd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("mustCwd:", err)
	}
	return wd
}

func selfUrlFunc(l []atom.Link) string {
	for _, l := range l {
		if l.Rel == "self" {
			return l.Href
		}
	}
	return "#"
}

func altUrlFunc(l []atom.Link) string {
	for _, l := range l {
		if l.Rel == "alternate" {
			return l.Href
		}
	}
	return "#"
}
