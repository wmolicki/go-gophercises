package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func getText(node *html.Node) string {
	// if element node, go inside this node

	var text string
	var f func(*html.Node)

	f = func(node *html.Node) {
		if node.Type == html.TextNode {
			text = text + node.Data
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}

	f(node)

	return strings.Join(strings.Fields(text), " ")
}

func getLinkFromANode(aNode *html.Node) Link {
	link := Link{Text: getText(aNode)}
	for _, attr := range aNode.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}
	return link
}

func getLinksFromRootNode(root *html.Node) []Link {
	result := []Link{}

	var f func(*html.Node)

	f = func(node *html.Node) {

		if node.Type == html.ElementNode && node.Data == "a" {
			result = append(result, getLinkFromANode(node))
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}

	f(root)

	return result
}

func Parse(reader io.Reader) []Link {
	root, err := html.Parse(reader)
	if err != nil {
		log.Fatalf("could not load html tree: %v", err)
	}

	links := getLinksFromRootNode(root)

	return links
}

func main() {
	htmlFilePtr := flag.String("f", "", "path to html file")
	flag.Parse()
	if *htmlFilePtr == "" {
		log.Fatalf("you have to provide path to html file")
	}

	reader, err := os.Open(*htmlFilePtr)
	if err != nil {
		log.Fatalf("could not open file %s: %v", *htmlFilePtr, err)
	}

	links := Parse(reader)

	for i, link := range links {
		log.Printf("%d: %s -> %s", i, link.Text, link.Href)
	}
}
