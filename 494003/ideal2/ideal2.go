package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Define a structure to represent the XML data
type Person struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}

func main() {
	// Parse XML file with error recovery and auto-correction
	data, err := parseXMLFile("data.xml")
	if err != nil {
		log.Printf("Error parsing XML file: %v", err)
		return
	}
	fmt.Println("Parsed Data:", data)

	// Create XML file with corrected structure
	err = createXMLFile("output.xml", data)
	if err != nil {
		log.Printf("Error creating XML file: %v", err)
	}
}

// Function to parse XML file and handle errors
func parseXMLFile(filename string) ([]Person, error) {
	var people []Person

	// Read XML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Attempt to auto-correct XML issues
	autoCorrectedData := autoCorrectMissingTags(data)

	// Parse the corrected XML
	err = xml.Unmarshal(autoCorrectedData, &people)
	if err != nil {
		log.Printf("Error parsing corrected XML: %v", err)
		return nil, fmt.Errorf("error parsing corrected XML: %v", err)
	}
	return people, nil
}

// Function to auto-correct common XML issues
func autoCorrectMissingTags(data []byte) []byte {
	xmlString := string(data)

	// Fix mismatched closing tags by balancing opening and closing tags
	openTags := []string{}
	correctedXML := strings.Builder{}

	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch token := tok.(type) {
		case xml.StartElement:
			openTags = append(openTags, token.Name.Local)
			correctedXML.WriteString(fmt.Sprintf("<%s>", token.Name.Local))
		case xml.EndElement:
			if len(openTags) > 0 && openTags[len(openTags)-1] == token.Name.Local {
				openTags = openTags[:len(openTags)-1] // Pop matching start tag
				correctedXML.WriteString(fmt.Sprintf("</%s>", token.Name.Local))
			} else {
				// Skip mismatched end tag
				log.Printf("Skipping mismatched closing tag: </%s>", token.Name.Local)
			}
		case xml.CharData:
			correctedXML.WriteString(string(token))
		}
	}

	// Add missing closing tags
	for i := len(openTags) - 1; i >= 0; i-- {
		correctedXML.WriteString(fmt.Sprintf("</%s>", openTags[i]))
	}

	return []byte(correctedXML.String())
}

// Function to create an XML file with corrected data
func createXMLFile(filename string, people []Person) error {
	// Marshal data with indentation for pretty printing
	output, err := xml.MarshalIndent(people, "", "  ")
	if err != nil {
		log.Printf("Error marshaling XML: %v", err)
		return fmt.Errorf("error marshaling XML: %v", err)
	}

	// Write XML to the file
	err = ioutil.WriteFile(filename, output, 0644)
	if err != nil {
		log.Printf("Error writing XML file: %v", err)
		return fmt.Errorf("error writing XML file: %v", err)
	}

	return nil
}
