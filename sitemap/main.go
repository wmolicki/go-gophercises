package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wmolicki/gophercises/linkparser"
)

func Build(baseUrl string) []string {

	result := []string{}
	visited := make(map[string]bool)

	var f func(string)

	f = func(u string) {

		u = normalizeUrl(baseUrl, u)

		body, err := fetch(u)
		if err != nil {
			log.Fatalf("could not fetch %s: %v", u, err)
		}

		links := linkparser.Parse(body)

		for _, link := range getLinksFromSameDomain(baseUrl, links) {
			if !visited[link.Href] {
				visited[link.Href] = true
				result = append(result, link.Href)
				f(link.Href)
			}
		}
	}

	f(baseUrl)

	return result
}

func getLinksFromSameDomain(base string, links []linkparser.Link) (out []linkparser.Link) {
	for _, link := range links {
		if sameDomain(base, link.Href) {
			out = append(out, link)
		}
	}
	return out
}

func normalizeUrl(base, u string) string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		log.Fatalf("could not parse %s: %v", base, err)
	}
	uUrl, err := url.Parse(u)
	if err != nil {
		log.Fatalf("could not parse %s: %v", u, err)
	}
	uUrl.Fragment = ""
	uUrl.RawQuery = ""
	if uUrl.Host == "" {
		uUrl.Host = baseUrl.Host
		uUrl.Scheme = baseUrl.Scheme
	}

	return uUrl.String()
}

func sameDomain(u1, u2 string) bool {
	// Returns whether u2 is in the same domain as u1.

	if u2[0] == '/' {
		// path links are always in the same domain
		return true
	}

	if u1[0] != '/' && !strings.HasPrefix(u1, "https://") && !strings.HasPrefix(u1, "http://") {
		u1 = "https://" + u1
	}

	url1, err := url.Parse(u1)
	if err != nil {
		log.Fatalf("could not parse %s: %v", u1, err)
	}

	url2, err := url.Parse(u2)
	if err != nil {
		log.Fatalf("could not parse %s: %v", u2, err)
	}

	return url1.Host == url2.Host
}

func fetch(url string) (io.Reader, error) {
	client := http.Client{Timeout: time.Duration(60) * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch %s: %v", url, err)
	}
	return r.Body, nil
}

func main() {
	fmt.Println("sitemap")
	link := flag.String("t", "", "target site to map")
	flag.Parse()
	if *link == "" {
		log.Fatalf("Target site flag required")
	}
	fmt.Printf("Preparing site map for %s\n", *link)

	for i, a := range Build(*link) {
		fmt.Printf("%d: %s", i, a)
	}
}
