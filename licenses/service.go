package licenses

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/meso-org/meso-license-service/repository"
	"golang.org/x/net/html"
)

type Service interface {
	StoreLicense(lic repository.License) (repository.License, error)
	UpdateLicense(lic repository.License) error
	VerifyLicense(lic repository.License) error
}

type service struct {
	licenses repository.LicenseRepository
}

func (s *service) StoreLicense(lic repository.License) (repository.License, error) {
	s.licenses.Store(&lic)
	return lic, nil
}

//updateLicense seems redundant. keeping for now
func (s *service) UpdateLicense(lic repository.License) error {
	s.licenses.Store(&lic)
	return nil
}

func (s *service) VerifyLicense(lic repository.License) error {
	createDcaPost(lic)
	return nil
}

//post to department of consumer affairs website
func createDcaPost(license repository.License) {
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
		log.Print("createDcaPost reading request:")
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		log.Print("createDcaPost:")
		log.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("createDcaPost:")
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
func htmlNodeTraversal(n *html.Node, license repository.License) {
	if n.Type == html.ElementNode && n.Data == "ul" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "actions" {
				text := &bytes.Buffer{}
				collectionBuffer := collectText(n, text)
				collectedText := collectionBuffer.String()
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

func verifyCollectedText(s string, license repository.License) {
	//we would need to pass this for user specific data
	/*
		name := "LASTNAME, FIRSTNAME"
		number := "633681"
		licenseType := "Registered Nurse"
	*/

	//Format Name: join first and last. then make them uppercase
	name := strings.ToUpper(license.LastName + ", " + license.FirstName)
	number := strconv.Itoa(license.Number)
	licenseType := license.LicenseDesc.Name

	matchExpression := name + "+\\s+License Number:+\\s+" + number + "+\\s+License Type:+\\s+" + licenseType
	//expression: return true if string matches FirstName LastName + License Name and Type + License Status

	match := expressionToRegex(matchExpression).MatchString(s)
	if match == true {
		newExpression := "[\n\r].*License Status:\\s*([^\n\r]*)"
		//returns string array [0] being "License Status: whateverstatus"

		result := expressionToRegex(newExpression).FindAllString(s, 1)
		if result == nil {
			log.Printf("verifyCollectedText: regex check nil")
		} else {
			license.Verify = true
			extractedResult := strings.Split(result[0], ":")
			license.Status = extractedResult[len(extractedResult)-1]
			license.Expiration = expirationDate(s)
			log.Println("Verified license: " + strconv.Itoa(license.Number))
		}
	} else {
		license.Verify = false
		log.Println("verifyCollectedText: license requested has no match")
	}
}

//Helper function for creating regex expressions
func expressionToRegex(expression string) *regexp.Regexp {
	var regex = regexp.MustCompile(expression)
	return regex
}

func expirationDate(s string) string {
	expression := "\\w+\\s\\d{2},\\s\\d{4}"
	index := expressionToRegex(expression).FindStringSubmatch(s)
	return index[0]
}

func NewService(licenseRepository repository.LicenseRepository) Service {
	return &service{
		licenses: licenseRepository,
	}
}
