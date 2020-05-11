package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/markbates/uad"

	"github.com/gobuffalo/flect"
	"github.com/gocolly/colly"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mm := map[string]uad.Plugin{}

	u := "https://www.uaudio.com/uad-plugins/all-plugins.html"
	root := "https://www.uaudio.com/uad-plugins/"

	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		h := e.Attr("href")
		if !strings.HasPrefix(h, root) || strings.Contains(h, "bundle") {
			return
		}

		p := mm[h]
		p.URL = h

		cat := strings.TrimPrefix(h, root)
		cat = strings.Split(cat, "/")[0]
		p.Category = flect.Titleize(cat)

		mm[h] = p
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		h := e.Text
		url := e.Request.URL.String()
		p := mm[url]
		p.Name = h
		p.URL = url
		mm[url] = p
	})

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		h := e.Text
		url := e.Request.URL.String()
		p := mm[url]
		if len(p.Description) == 0 {
			p.Description = fmt.Sprintf("%s\n", h)
		}
		p.URL = url
		mm[url] = p
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(u)

	f, err := os.Create("plugins.json")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(mm); err != nil {
		return err
	}
	return nil
}
