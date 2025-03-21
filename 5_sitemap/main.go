package main

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ranmerc/gophercises/link/parser"
)

type SitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func main() {
	url := "https://www.calhoun.io"

	links, err := getPageLink(url)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(0)
	}

	links = filterMap(links, func(link parser.Link) bool {
		return strings.HasPrefix(link.Href, "/") || strings.HasPrefix(link.Href, url)
	}, func(link parser.Link) parser.Link {
		if strings.HasPrefix(link.Href, "/") {
			return parser.Link{
				Text: link.Text,
				Href: url + link.Href,
			}
		}

		return link
	})

	sitemapURLs := []SitemapURL{}

	seen := make(map[string]bool)

	seen[url] = true
	seen["/"] = true

	i := 0
	for i < len(links) {
		link := links[i].Href

		if ok := seen[link]; ok {
			i += 1
			continue
		}

		sitemapURLs = append(sitemapURLs, SitemapURL{
			Loc: link,
		})

		seen[link] = true

		linksOnPage, err := getPageLink(link)
		if err != nil {
			slog.Error(err.Error())
			i += 1
			continue
		}

		for _, link := range linksOnPage {
			href := link.Href
			if strings.HasPrefix(href, url) {
				links = append(links, link)
			}

			if strings.HasPrefix(href, "/") {
				links = append(links, parser.Link{
					Text: link.Text,
					Href: url + href,
				})
			}
		}

		i += 1
	}

	xmlBytes, err := xml.MarshalIndent(sitemapURLs, " ", " ")

	fmt.Println(xml.Header + string(xmlBytes))
}

func getPageLink(url string) ([]parser.Link, error) {
	slog.Info("Visiting " + url)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to GET page %v", err)
	}

	links, err := parser.ParseLinks(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse page content for links %v", err)
	}

	return links, nil
}

func filterMap[T any](ss []T, test func(T) bool, transform func(T) T) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, transform(s))
		}
	}
	return
}
