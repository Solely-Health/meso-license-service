package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	// "strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/gorilla/mux"
)

type LicenseType struct {
	boardCode   string
	name string 
	licenseCode string
}

type Status struct {
	current  bool 
	delinquent bool 
	deceased bool
	voluntarySurrender bool
}
type License struct {
	number      int
	licenseType LicenseType
	status Status 
	expiration string 
	description string
	secondaryStatus string
}

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
	//payload := strings.NewReader("boardCode=0&licenseType=0&firstName=RUBY&lastName=ABRANTES&licenseNumber=633681")
	payload := strings.NewReader("boardCode=0&licenseType=224&firstName=RUBY&lastName=ABRANTES&licenseNumber=633681")

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
				collectionBuffer := collectText(n, text)
				collectedText := collectionBuffer.String()
				//fmt.Printf(collectedText)
				verifyCollectedText(collectedText)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		htmlNodeTraversal(c)
	}
}

//Go through the tree and write to buffer
func collectText(n *html.Node, buf *bytes.Buffer) *bytes.Buffer {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
	return buf
}

func verifyCollectedText(s string) {
	//hardcoded we would need to change this for user specific data
	name := "ABRANTES, RUBY" 
	number := "633681"
	licenseType := "Registered Nurse"

	expression := name + "+\\s+License Number:+\\s+" + number + "+\\s+License Type:+\\s+" + licenseType 
	// "+\\s+License Status: Current"
	//expression says return true if FirstName LastName + License Name and Type + License Status == Current.

	var validID = regexp.MustCompile(expression)

	match := validID.MatchString(s)

	if match == true {
		expression := "License Status: Current"
		var verifyStatus = regexp.MustCompile(expression)
		if verifyStatus.MatchString(s){
			//return true
			fmt.Printf("True")
		}else {
			//return false
			fmt.Printf("False")
		}
	} else {
		//return false
		fmt.Printf("False")
	}


	//fmt.Println("Match Result: " + name + " " + number + " " + licenseType + ": " + strconv.FormatBool(validID.MatchString(s)))

}

func main() {
	fmt.Printf("Started Service\n")
	//manually doing a post. for testing
	createDcaPost()
 
	//handler not setup
	router := mux.NewRouter()
	router.HandleFunc("/license", licenseRequest)
	//log.Fatal(http.ListenAndServe(":8080", router))
}
