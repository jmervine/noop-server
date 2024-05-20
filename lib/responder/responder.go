package responder

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Responders map[string]Responder

func (r *Responders) Match(key string) (Responder, bool) {
	if responder, ok := (*r)[key]; ok {
		return responder, true
	}

	for endpoint, responder := range *r {
		if endpoint == "*" {
			return responder, true
		}

		match, err := regexp.MatchString(endpoint, key)
		if err != nil {
			log.Fatal(err)
		}

		if match {
			return responder, true
		}
	}

	return Responder{}, false
}

type Responder struct {
	Status   int
	Sleep    uint
	Response string
}

func Load(filename string) (*Responders, error) {
	responders := new(Responders)

	filename, err := filepath.Abs(filename)
	if err != nil {
		return responders, err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return responders, err
	}

	err = yaml.Unmarshal(content, &responders)
	return responders, err
}
