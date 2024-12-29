package main

import (
	"fmt"
)

type Person struct {
	Name   string
	Age    int
	Height float64
}

func main() {
	person := Person{Name: "Alice", Age: 30, Height: 5.9}
	cities := []string{"New York", "Los Angeles", "Chicago"}
	visits := map[string]int{"New York": 5, "Los Angeles": 3, "Chicago": 2}

	formattedString := fmt.Sprintf("Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
		person.Name, person.Age, person.Height, len(cities), cities,
		visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])

	fmt.Println(formattedString)
}