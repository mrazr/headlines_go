package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func IsHeading(e atom.Atom) bool {
	return e == atom.H1 || e == atom.H2 || e == atom.H3 || e == atom.H4 || e == atom.H5 || e == atom.H6
}

func ExtractText(n *html.Node) (string, error) {
	ii := n.Descendants()
	var text string = ""
	ii(func(node *html.Node) bool {
		if node.Type != html.TextNode {
			return true
		}
		text = node.Data
		return false
	})
	return text, nil
}

func ExtractHeaders(content io.Reader) (map[int][]string, error) {
	doc, err := html.Parse(content)
	if err != nil {
		return nil, err
	}
	h_el_to_index := map[atom.Atom]int{
		atom.H1: 0,
		atom.H2: 1,
		atom.H3: 2,
		atom.H4: 3,
		atom.H5: 4,
		atom.H6: 5,
	}
	headers_map := make(map[int][]string)
	for i := 0; i < 6; i += 1 {
		headers_map[i] = make([]string, 0, 20)
	}
	el_iter := doc.Descendants()
	el_iter(func(node *html.Node) bool {
		if node.Type == html.ElementNode && IsHeading(node.DataAtom) {
			text, _ := ExtractText(node)
			if text != "" {
				headers_map[h_el_to_index[node.DataAtom]] = append(headers_map[h_el_to_index[node.DataAtom]], text)
			}
		}
		return true
	})
	return headers_map, nil
}

func main() {
	file_contents, err := os.ReadFile("./sites.txt")

	if err != nil {
		fmt.Println("Could not read the file")
		os.Exit(1)
	}
	text := string(file_contents)
	websites := strings.Split(text, "\n")
	var headlines map[int][]string
	for _, website := range websites {
		fmt.Printf("got %s\n", website)
		resp, err := http.Get(website)
		if err != nil {
			fmt.Printf("Could not GET %s\n", website)
		} else {
			headlines, _ = ExtractHeaders(resp.Body)
		}
	}
	for _, headline := range headlines[2] {
		fmt.Println(headline)
	}
}
