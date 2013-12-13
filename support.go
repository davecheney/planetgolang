package main

import (
	"log"
	"os"
)

func mustCwd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("mustCwd:", err)
	}
	return wd
}
