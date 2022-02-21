package feeds

import (
	"encoding/xml"
	"strings"
)

const (
	GeneratorWordpress = 1
)

type GeneratorChannel struct {
	Generator string `xml:"generator"`
}

type GeneratorFeed struct {
	Channel MediumChannel `xml:"channel"`
}

func FindGenerator(url string) (int, []byte, error) {

	generators := map[string]int{
		"wordpress": GeneratorWordpress,
	}

	bod, err := httpget(url)
	if err != nil {
		return 0, bod, err
	}
	var f GeneratorFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		return 0, bod, nil // this is ok actually... just return 0 for type
	}

	gen := 0
	for key, element := range generators {
		if strings.Contains(f.Channel.Generator, key) {
			gen = element
		}
	}

	return gen, bod, nil
}
