package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
)

// Define a structure to represent the XML data
type Person struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}

func main() {
	// Parse XML file with error recovery
	data, err := parseXMLFile("data.xml")
	if err != nil {
		log.Printf("Error parsing XML file: %v", err)
		return
	}
	fmt.Println("Parsed Data:", data)

	// Create XML file with error recovery
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

	// Use xml.Unmarshal to parse the XML data into the People struct
	err = xml.Unmarshal(data, &people)
	if err != nil {
		// Handle specific parsing errors
		if syntaxError, ok := err.(*xml.SyntaxError); ok {
			log.Printf("Malformed XML syntax: %v", syntaxError)
			return nil, fmt.Errorf("XML syntax error: %v", syntaxError)
		}
		// Handle other errors
		log.Printf("Error parsing XML: %v", err)
		return nil, fmt.Errorf("error parsing XML: %v", err)
	}

	return people, nil
}

// Function to create XML file and handle errors
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
