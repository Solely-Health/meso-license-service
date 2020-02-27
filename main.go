package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/net/html"
)

type LicenseType struct {
	BoardCode   int    `json:"boardCode"`
	Name        string `json:"licenseName"`
	LicenseCode int    `json:"licenseCode"`
}

type License struct {
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Number          int         `json:"licenseNumber"`
	LicenseDesc     LicenseType `json:"licenseType"`
	Status          string
	Expiration      string
	Description     string
	SecondaryStatus string
	Verify          bool
}

func licenseRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got the request")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Fprintf(w, "It broke")
	}
	//fmt.Println(body)
	var newLicense License
	if err := json.Unmarshal(body, &newLicense); err != nil {
		log.Println(err)
	}
	fmt.Println("printing License object")
	fmt.Println(newLicense)
	createDcaPost(&newLicense)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLicense)
}

//post to department of consumer affairs website
func createDcaPost(license *License) {
	//Hardcoded payload example.
	//payload := strings.NewReader("boardCode=0&licenseType=224&firstName=RUBY&lastName=ABRANTES&licenseNumber=633681")

	url := "https://search.dca.ca.gov/results"
	method := "POST"
	board := strconv.Itoa(license.LicenseDesc.BoardCode)
	licenseCode := strconv.Itoa(license.LicenseDesc.LicenseCode)
	licenseNumber := strconv.Itoa(license.Number)
	firstName := license.FirstName
	lastName := license.LastName

	//create payload for POST
	payload := strings.NewReader("boardCode=" + board + "&licenseType=" + licenseCode + "&firstName=" + firstName + "&lastName=" + lastName + "&licenseNumber=" + licenseNumber)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	//takes string of html. parses to nodes. (basically makes a tree of tags)
	parseMe, err := html.Parse(strings.NewReader(string(body)))
	htmlNodeTraversal(parseMe, license)
	if err != nil {
		log.Fatal(err)
	}
}

//Finds tag we need and collects into text
func htmlNodeTraversal(n *html.Node, license *License) {
	if n.Type == html.ElementNode && n.Data == "ul" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "actions" {
				text := &bytes.Buffer{}
				collectionBuffer := collectText(n, text)
				collectedText := collectionBuffer.String()
				//fmt.Printf(collectedText)
				verifyCollectedText(collectedText, license)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		htmlNodeTraversal(c, license)
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

func verifyCollectedText(s string, license *License) {
	//we would need to pass this for user specific data
	/*
		name := "ABRANTES, RUBY"
		number := "633681"
		licenseType := "Registered Nurse"
	*/

	//Format Name: join first and last. then make them uppercase
	name := strings.ToUpper(license.LastName + ", " + license.FirstName)
	number := strconv.Itoa(license.Number)
	licenseType := license.LicenseDesc.Name

	matchExpression := name + "+\\s+License Number:+\\s+" + number + "+\\s+License Type:+\\s+" + licenseType
	// "+\\s+License Status: Current"
	//expression says return true if FirstName LastName + License Name and Type + License Status == Current.
	match := expressionToRegex(matchExpression).MatchString(s)
	if match == true {
		newExpression := "[\n\r].*License Status:\\s*([^\n\r]*)"
		//returns string array [0] being "License Status: whateverstatus"
		result := expressionToRegex(newExpression).FindAllString(s, 1)
		if result == nil {
			fmt.Print("we broke")
		} else {
			extractedResult := strings.Split(result[0], ":")
			license.Status = extractedResult[len(extractedResult)-1]
			license.Expiration = expirationDate(s)
			fillLicenseObject(s, license)
		}
	} else {
		fillLicenseObject(s, license)
		fmt.Printf("False! no match")
		//fmt.Printf(s)
	}
}

//Helper function for creating regex expressions
func expressionToRegex(expression string) *regexp.Regexp {
	var regex = regexp.MustCompile(expression)
	return regex
}

func fillLicenseObject(s string, license *License) {
	//set status of license sent from verifyCollectedText()
	license.Verify = true
	fmt.Println("made it to fill function")
	fmt.Println(license)
}

func expirationDate(s string) string {
	expression := "\\w+\\s\\d{2},\\s\\d{4}"
	index := expressionToRegex(expression).FindStringSubmatch(s)
	return index[0]
}

func main() {
	fmt.Printf("Started Service\n")
	//manually doing a post. for testing
	/*
		var license *License
		license = new(License)
		createDcaPost(license)
		fmt.Println(license.status)
	*/

	router := mux.NewRouter()
	router.HandleFunc("/license", licenseRequest)

	//for local testing
	log.Fatal(http.ListenAndServe("localhost:8080", router))

	//log.Fatal(http.ListenAndServe(":8080", router))
}
