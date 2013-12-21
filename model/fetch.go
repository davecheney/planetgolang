package model

import (
	"encoding/xml"
	"errors"
	"net/http"

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
