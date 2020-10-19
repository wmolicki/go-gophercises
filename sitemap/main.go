package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/wmolicki/gophercises/linkparser"
)

func main() {
	fmt.Println("sitemap")
	link := flag.String("t", "", "target site to map")
	flag.Parse()
	if *link == "" {
		log.Fatalf("Target site flag required")
	}
	fmt.Printf("Preparing site map for %s\n", *link)

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	r, err := client.Get(*link)
	if err != nil {
		log.Fatalf("could not fetch %s: %v", *link, err)
	}

	links := linkparser.Parse(r.Body)
	for _, link := range links {
		fmt.Printf("%s -> %s\n", link.Text, link.Href)
	}
}
