package main

import (
	"io/ioutil"
	"log"
	"os"

	"encoding/json"

	"github.com/google/go-github/v24/github"
)

func loadJSONFile(path string) *github.PushEvent {
	jsonFile, err := os.Open(path)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var data github.PushEvent

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatal(err)
		return nil
	}

	return &data
}
