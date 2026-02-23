package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

var visited = make(map[string]bool)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go https://example.com")
		return
	}

	rawURL := os.Args[1]

	// start crawling with depth limit
	crawl(rawURL, 2)
}

func crawl(rawURL string, depth int) {
	if depth <= 0 {
		return
	}

	// prevent duplicate crawling
	if visited[rawURL] {
		return
	}
	visited[rawURL] = true

	fmt.Println("Crawling:", rawURL)

	resp, err := http.Get(rawURL)
	if err != nil {
		log.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Bad status:", resp.Status)
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("Parse error:", err)
		return
	}

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	links := extractLinks(doc, baseURL)

	for _, link := range links {
		crawl(link, depth-1)
	}
}

func extractLinks(n *html.Node, base *url.URL) []string {
	var links []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {

				link, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}

				absolute := base.ResolveReference(link)

				// skip different domains
				if absolute.Host != base.Host {
					continue
				}

				// skip fragments
				if absolute.Fragment != "" {
					continue
				}

				links = append(links, absolute.String())
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childLinks := extractLinks(c, base)
		links = append(links, childLinks...)
	}

	return links
}
