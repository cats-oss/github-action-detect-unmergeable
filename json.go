package main

import (
	"io/ioutil"
	"log"
	"os"

	"encoding/json"
)

// PushEventData is workaround to use *github.PushEvent.
// If we available, we should use *github.PushEvent.
// But we cannot use it by parsing error. See https://github.com/cats-oss/github-action-detect-unmergeable/issues/71
type PushEventData struct {
	Ref     *string `json:"ref,omitempty"`
	Compare *string `json:"compare,omitempty"`
}

func (p *PushEventData) GetRef() string {
	if p == nil || p.Ref == nil {
		return ""
	}
	return *p.Ref
}

func (p *PushEventData) GetCompare() string {
	if p == nil || p.Compare == nil {
		return ""
	}
	return *p.Compare
}

func loadJSONFile(path string) *PushEventData {
	jsonFile, err := os.Open(path)
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var data PushEventData

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatal(err)
		return nil
	}

	return &data
}
