package main

import (
	"os"
	"log"
)

func mustCwd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("mustCwd:",err)
	}
	return wd
}
