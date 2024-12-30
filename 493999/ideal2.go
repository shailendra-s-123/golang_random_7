
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Params struct for managing query parameters
type Params struct {
	data map[string]interface{}
}

// NewParams initializes a new Params instance
func NewParams() *Params {
	return &Params{
		data: make(map[string]interface{}),
	}
}

// Set stores a value for a specific key
func (p *Params) Set(key string, value interface{}) {
	p.data[key] = value
}

// Get retrieves a value for a specific key
func (p *Params) Get(key string) (interface{}, bool) {
	v, ok := p.data[key]
	return v, ok
}

// String retrieves a string value for a specific key
func (p *Params) String(key string) (string, error) {
	v, ok := p.Get(key)
	if !ok {
		return "", fmt.Errorf("parameter '%s' not found", key)
	}
	return fmt.Sprintf("%v", v), nil
}

// Int retrieves an integer value for a specific key
func (p *Params) Int(key string) (int, error) {
	v, ok := p.Get(key)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("parameter '%s' expected int, got %T", key, v)
	}
}

// Bool retrieves a boolean value for a specific key
func (p *Params) Bool(key string) (bool, error) {
	v, ok := p.Get(key)
	if !ok {
		return false, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case string:
		return strconv.ParseBool(v)
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("parameter '%s' expected bool, got %T", key, v)
	}
}

// Float64 retrieves a float64 value for a specific key
func (p *Params) Float64(key string) (float64, error) {
	v, ok := p.Get(key)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("parameter '%s' expected float64, got %T", key, v)
	}
}

// Date retrieves a date value for a specific key, supporting a format
func (p *Params) Date(key, format string) (time.Time, error) {
	v, ok := p.Get(key)
	if !ok {
		return time.Time{}, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case string:
		return time.Parse(format, v)
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("parameter '%s' expected date, got %T", key, v)
	}
}

// Slice retrieves a slice of strings for a specific key
func (p *Params) Slice(key string) ([]string, error) {
	v, ok := p.Get(key)
	if !ok {
		return nil, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case []string:
		return v, nil
	default:
		return nil, fmt.Errorf("parameter '%s' expected slice of strings, got %T", key, v)
	}
}

// Map retrieves a map of strings for a specific key
func (p *Params) Map(key string) (map[string]string, error) {
	v, ok := p.Get(key)
	if !ok {
		return nil, fmt.Errorf("parameter '%s' not found", key)
	}
	switch v := v.(type) {
	case map[string]string:
		return v, nil
	default:
		return nil, fmt.Errorf("parameter '%s' expected map of strings, got %T", key, v)
	}
}

// ParseQuery parses a URL query string into Params struct
func (p *Params) ParseQuery(query string) error {
	values, err := url.ParseQuery(query)
	if err != nil {
		return err
	}
	for key, vals := range values {
		// Handle repeated parameters by storing as slice
		if len(vals) > 1 {
			p.Set(key, vals)
		} else {
			p.Set(key, vals[0])
		}
	}
	return nil
}

// AggregateInts performs aggregation on a set of integer values for a specific key
func (p *Params) AggregateInts(key string, operation func(int, int) int) (int, error) {
	vals, ok := p.data[key].([]string)
	if !ok {
		return 0, fmt.Errorf("parameter '%s' expected []string, got %T", key, p.data[key])
	}
	var result int
	for _, v := range vals {
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid integer '%s' for parameter '%s'", v, key)
		}
		result = operation(result, i)
	}
	return result, nil
}

// SerializeJSON serializes the Params struct to JSON
func (p *Params) SerializeJSON() ([]byte, error) {
	return json.Marshal(p.data)
}

// DeserializeJSON deserializes JSON into the Params struct
func (p *Params) DeserializeJSON(data []byte) error {
	return json.Unmarshal(data, &p.data)
}

// Example demonstrating parsing, validation, aggregation, and serialization
func main() {
	params := NewParams()

	// Example query string
	query := "name=Alice&age=30&is_active=true&tags=sport,music&dates=2023-10-01,2023-10-02&pagination=2&sort=asc"
	err := params.ParseQuery(query)
	if err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	// Extracting values
	name, _ := params.String("name")
	age, _ := params.Int("age")
	isActive, _ := params.Bool("is_active")
	dates, _ := params.Slice("dates")
	tags, _ := params.Slice("tags")
	pagination, _ := params.Int("pagination")
	sort, _ := params.String("sort")

	// Display parsed values
	fmt.Println("Parsed Values:")
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("Is Active:", isActive)
	fmt.Println("Dates:", dates)
	fmt.Println("Tags:", tags)
	fmt.Println("Pagination:", pagination)
	fmt.Println("Sort:", sort)

	// Aggregation example
	tagCount, _ := params.AggregateInts("tags", func(a, b int) int {
		return a + 1
	})
	fmt.Println("Number of Tags:", tagCount)

	// Serialization and Deserialization example
	serializedData, _ := params.SerializeJSON()
	fmt.Println("Serialized Data:", string(serializedData))

	// Deserializing back into Params
	newParams := NewParams()
	err = newParams.DeserializeJSON(serializedData)
	if err != nil {
		fmt.Println("Error deserializing:", err)
	} else {
		fmt.Println("Deserialized Data:", newParams.data)
	}
}

