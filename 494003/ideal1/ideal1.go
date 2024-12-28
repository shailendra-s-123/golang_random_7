package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Structs for XML Data
type Person struct {
	Name    string `xml:"name" validate:"required"`
	Age     int    `xml:"age" validate:"min=0"`
	Address string `xml:"address" validate:"required"`
}

type Data struct {
	XMLName xml.Name `xml:"data"`
	People  []Person `xml:"person"`
}

// Custom Error Types
var (
	ErrMalformedXML = errors.New("malformed XML data")
	ErrMissingField = errors.New("missing required XML fields")
)

// Read XML file
func readXML(filePath string) (Data, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Data{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	var data Data
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Data{}, fmt.Errorf("error reading XML token: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "data" {
				if err := decoder.DecodeElement(&data, &t); err != nil {
					return Data{}, fmt.Errorf("failed to decode XML element: %w", err)
				}
			}
		}
	}

	// Check for missing required fields
	for _, person := range data.People {
		if strings.TrimSpace(person.Name) == "" || strings.TrimSpace(person.Address) == "" {
			return Data{}, fmt.Errorf("%w: name or address missing in person element", ErrMissingField)
		}
	}
	return data, nil
}

// Write XML file
func writeXML(filePath string, data Data) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding XML: %w", err)
	}
	return nil
}

func main() {
	inputFile := "data.xml"
	outputFile := "output.xml"

	log.Println("Reading XML file...")
	data, err := readXML(inputFile)
	if err != nil {
		if errors.Is(err, ErrMalformedXML) || errors.Is(err, ErrMissingField) {
			log.Fatalf("Critical error with XML structure: %v", err)
		}
		log.Fatalf("Error reading XML: %v", err)
	}

	log.Printf("Successfully read XML: %+v\n", data)

	// Modify data for demonstration
	for i := range data.People {
		data.People[i].Age += 1 // Increment age for all people
	}

	log.Println("Writing updated XML file...")
	if err := writeXML(outputFile, data); err != nil {
		log.Fatalf("Error writing XML: %v", err)
	}

	log.Println("XML processing completed successfully!")
}
