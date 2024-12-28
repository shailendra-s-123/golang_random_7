package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/go-playground/validator.v9"
)

type Person struct {
	Name    string `xml:"name,attr" validate:"required"`
	Age     int    `xml:"age,attr" validate:"min=0"`
	Address string `xml:"address" validate:"required"`
}

type Data struct {
	XMLName xml.Name  `xml:"data"`
	People  []Person `xml:"person"`
}

func main() {
	// Validator for XML data
	validate := validator.New()

	// Read XML file
	xmlFile, err := os.Open("data.xml")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer xmlFile.Close()

	// Read the entire file content
	content, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Parse XML data
	var data Data
	err = xml.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("Error parsing XML: %v", err)
	}

	// Validate XML data
	if err := validate.Struct(data); err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	// Write XML file
	outFile, err := os.Create("output.xml")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer outFile.Close()

	// Encode XML data
	enc := xml.NewEncoder(outFile)
	enc.Indent("  ", "    ")
	if err := enc.Encode(data); err != nil {
		log.Fatalf("Error encoding XML: %v", err)
	}

	fmt.Println("XML file read, validated, and written successfully.")
}