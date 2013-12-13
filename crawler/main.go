package main

import (
	"net/http"
	"os"
	"log"
	"encoding/xml"

	"code.google.com/p/rsc/blog/atom"
)

func main() {
	resp, err := http.Get(os.Args[1])
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
	log.Println(feed)
}	
