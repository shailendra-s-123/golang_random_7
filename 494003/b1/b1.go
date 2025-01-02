package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "os"
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

func parseXMLFile(filename string) ([]Person, error) {
    var people []Person

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Use xml.Unmarshal to parse the XML data into the People struct
    err = xml.Unmarshal(data, &people)
    if err != nil {
        // Check for specific parsing errors like malformed tags
        if syntaxError, ok := err.(*xml.SyntaxError); ok {
            log.Printf("Malformed XML syntax at byte %d: %v", syntaxError.Offset, err)
            return nil, err
        }
        // Handle other parsing errors
        log.Printf("Error parsing XML: %v", err)
        return nil, err
    }

    return people, nil
}

func createXMLFile(filename string, people []Person) error {
    output, err := xml.MarshalIndent(people, "", "  ")
    if err != nil {
        log.Printf("Error marshaling XML: %v", err)
        return err
    }

    err = ioutil.WriteFile(filename, output, 0644)
    if err != nil {
        log.Printf("Error writing XML file: %v", err)
        return err
    }

    return nil
} 