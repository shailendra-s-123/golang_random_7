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

// Structs for XML Data with namespaces
type Person struct {
	XMLName  xml.Name `xml:"person"`
	Name     string   `xml:"ns1:name" validate:"required"`
	Age      int      `xml:"ns1:age" validate:"min=0"`
	Address string   `xml:"ns1:address" validate:"required"`
}

type Data struct {
	XMLName xml.Name `xml:"data"`
	People  []Person `xml:"person"`
}

// Custom Error Types
var (
	ErrMalformedXML    = errors.New("malformed XML data")
	ErrMissingField    = errors.New("missing required XML field")
	ErrMissingNamespace = errors.New("missing required XML namespace")
)

// Read XML file with namespace support
func readXMLWithNamespace(filePath string) (Data, error) {
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

	// Check for missing required fields and namespaces
	for _, person := range data.People {
		if strings.TrimSpace(person.Name) == "" || strings.TrimSpace(person.Address) == "" {
			return Data{}, fmt.Errorf("%w: name or address missing in person element", ErrMissingField)
		}
		// Check if the required namespace "ns1" is present
		if person.XMLName.Space != "ns1" {
			return Data{}, fmt.Errorf("%w: missing required namespace 'ns1' for element 'person'", ErrMissingNamespace)
		}
	}
	return data, nil
}

// Write XML file with namespace support
func writeXMLWithNamespace(filePath string, data Data) error {
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
	inputFile := "data_with_namespace.xml"
	outputFile := "output_with_namespace.xml"
	
	log.Println("Reading XML file with namespace support...")
	data, err := readXMLWithNamespace(inputFile)
	if err != nil {
		if errors.Is(err, ErrMalformedXML) || errors.Is(err, ErrMissingField) || errors.Is(err, ErrMissingNamespace) {
			log.Fatalf("Critical error with XML structure: %v", err)
		}
		log.Fatalf("Error reading XML: %v", err)
	}
	log.Printf("Successfully read XML: %+v\n", data)

	// Modify data for demonstration
	for i := range data.People {
		data.People[i].Age += 1 // Increment age for all people
	}
	log.Println("Writing updated XML file with namespace support...")
}