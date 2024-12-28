package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Sample XML struct definition
type Item struct {
	Name    string `xml:"name"`
	Quantity int    `xml:"quantity"`
}

// Function to read an XML file
func readXMLFile(filePath string) ([]Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var items []Item
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading XML token: %v", err)
			return nil, err
		}
		if startElement, ok := token.(xml.StartElement); ok && startElement.Name.Local == "item" {
			var item Item
			if err := decoder.Decode(&item); err != nil {
				log.Printf("Error decoding item: %v", err)
				return nil, err
			}
			items = append(items, item)
		}
	}
	return items, nil
}

// Function to write XML data to a file
func writeXMLFile(filePath string, items []Item) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	err = encoder.StartElement(xml.StartElement{Name: xml.Name{Local: "items"}})
	if err != nil {
		log.Printf("Error starting XML encoding: %v", err)
		return err
	}

	for _, item := range items {
		err = encoder.Encode(item)
		if err != nil {
			log.Printf("Error encoding item: %v", err)
			return err
		}
	}

	err = encoder.EndElement(xml.EndElement{Name: xml.Name{Local: "items"}})
	if err != nil {
		log.Printf("Error ending XML encoding: %v", err)
		return err
	}

	return nil
}

func main() {
	const filePath = "items.xml"

	// Read XML data
	items, err := readXMLFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read XML file: %v", err)
	}

	fmt.Printf("Read items from file: %#v\n", items)

	// Modify items (for demonstration purposes)
	for i, item := range items {
		item.Quantity *= 2
		items[i] = item
	}

	// Write modified XML data
	if err := writeXMLFile(filePath, items); err != nil {
		log.Fatalf("Failed to write XML file: %v", err)
	}

	fmt.Printf("Modified items written to file.\n")
}