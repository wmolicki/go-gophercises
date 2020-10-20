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

		for _, url := range getUrlsFromSameDomain(baseUrl, links) {
			if !visited[url] {
				visited[url] = true
				result = append(result, url)
				f(url)
			}
		}
	}

	f(baseUrl)

	return result
}

func getUrlsFromSameDomain(base string, links []linkparser.Link) (out []string) {
	for _, link := range links {
		if fi := strings.IndexAny(link.Href, "#"); fi != -1 {
			link.Href = link.Href[:fi]
		}
		if link.Href != "" && sameDomain(base, link.Href) {
			out = append(out, link.Href)
		}
	}
	return out
}

func normalizeUrl(base, u string) string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		log.Fatalf("could not parse base %s: %v", base, err)
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

	if u2[0] == '/' || strings.IndexAny(u2, ".") == -1 {
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
		fmt.Printf("%d: %s\n", i+1, a)
	}
}
