package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/gorilla/mux"
)

func licenseRequest(w http.ResponseWriter, r *http.Request) {
	//for testing with docker
	fmt.Printf("got the request")

	//TODO take user input and format it in x-www-form-urlencoded
	//pass into dcaPost
	createDcaPost()
}

//post to department of consumer affairs website
func createDcaPost() {
	url := "https://search.dca.ca.gov/results"
	method := "POST"

	//Hardcoded an example nurse. Change when passing in user input
	payload := strings.NewReader("boardCode=0&licenseType=0&firstName=RUBY&lastName=ABRANTES&licenseNumber=633681")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	//takes string of html. parses to nodes. (basically makes a tree of tags)
	parseMe, err := html.Parse(strings.NewReader(string(body)))
	htmlNodeTraversal(parseMe)
	if err != nil {
		log.Fatal(err)
	}
}

//Finds tag we need and collects into text
func htmlNodeTraversal(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "ul" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "actions" {
				text := &bytes.Buffer{}
				collectText(n, text)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		htmlNodeTraversal(c)
	}
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
		fmt.Printf(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func main() {
	fmt.Printf("Started Service")

	//manually doing a post. for testing
	createDcaPost()

	//handler not setup
	router := mux.NewRouter()
	router.HandleFunc("/license", licenseRequest)
	log.Fatal(http.ListenAndServe(":8080", router))
}
