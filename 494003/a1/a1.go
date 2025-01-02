package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

// Configuration structure for your XML
type Config struct {
	Server string `xml:"server"`
	Port   int    `xml:"port"`
	Enable bool   `xml:"enable,attr"`
}

func parseXML(filePath string) (*Config, error) {
	config := Config{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("XML parsing error: %v", err)
			continue
		}
		switch se := tok.(type) {
		case xml.StartElement:
			err = decoder.StartElement(&se)
			if err != nil {
				log.Printf("XML parsing error: %v", err)
				continue
			}
		}
	}
	err = decoder.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing decoder: %v", err)
	}
	return &config, nil
}

func createXML(config *Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("  ", "    ") // Indent for pretty printing
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("error encoding XML: %v", err)
	}
	return nil
}

func main() {
	filePath := "config.xml"

	// Parse XML file
	config, err := parseXML(filePath)
	if err != nil {
		log.Fatalf("Failed to parse XML file: %v", err)
	}

	log.Printf("Parsed config: Server=%s, Port=%d, Enable=%v", config.Server, config.Port, config.Enable)

	// Create new XML file
	err = createXML(&Config{Server: "localhost", Port: 8080, Enable: true}, "output.xml")
	if err != nil {
		log.Fatalf("Failed to create XML file: %v", err)
	}

	log.Println("XML file created successfully.")
}