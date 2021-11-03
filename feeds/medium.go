package feeds

import (
	"encoding/xml"
	"fmt"
)

type MediumPost struct {
	Title string `xml:"title"`
	Desc  string `xml:"description"`
	Link  string `xml:"link"`
	Guid  string `xml:"guid"`
}

type MediumChannel struct {
	Title         string       `xml:"title"`
	Link          string       `xml:"link"`
	Desc          string       `xml:"description"`
	Image         string       `xml:"image"`
	Generator     string       `xml:"generator"`
	LastBuildData string       `xml:"lastBuildDate"`
	Items         []MediumPost `xml:"item"`
}

type MediumFeed struct {
	Channel MediumChannel `xml:"channel"`
}

func ParseMediumFeed(url string) {
	bod, err := httpget(url)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	var f MediumFeed
	if err := xml.Unmarshal(bod, &f); err != nil {
		fmt.Println("PARSE ERROR", err)
	}
	fmt.Printf("%+v\n", f)
}
