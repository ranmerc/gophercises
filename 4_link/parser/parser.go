package parser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link is the parsed html link.
type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("[Text: %q Href: %q]", l.Text, l.Href)
}

// ParseLinks parses out and returns the links in a HTML file.
func ParseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := []Link{}

	for node := range doc.Descendants() {
		if node.Type == html.ElementNode && node.Data == "a" {
			link := Link{}

			// Find the "href" attribute
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					link.Href = attr.Val
					break
				}
			}

			texts := []string{}
			for childNode := range node.Descendants() {
				if childNode.Type == html.TextNode {
					if len(childNode.Data) > 0 {
						texts = append(texts, childNode.Data)
					}
				}
			}
			link.Text = strings.TrimSpace(strings.Join(texts, ""))

			links = append(links, link)
		}
	}

	return links, nil
}
