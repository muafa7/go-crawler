package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go https://example.com")
		return
	}

	rawURL := os.Args[1]

	resp, err := http.Get(rawURL)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad Status: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	parsedURL, _ := url.Parse(rawURL)

	extractLinks(doc, parsedURL)
}

func extractLinks(n *html.Node, base *url.URL) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}

				absolute := base.ResolveReference(link)

				if absolute.Host == base.Host {
					fmt.Println(absolute.String())
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, base)
	}
}
