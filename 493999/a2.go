package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	data map[string]interface{}
}

func NewParams() *Params {
	return &Params{
		data: make(map[string]interface{}),
	}
}

func (p *Params) Set(key string, value interface{}) {
	p.data[key] = value
}

func (p *Params) Get(key string) interface{} {
	return p.data[key]
}

func (p *Params) String(key string) string {
	v := p.Get(key)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (p *Params) Int(key string) (int, error) {
	v := p.Get(key)
	switch s := v.(type) {
	case string:
		return strconv.Atoi(s)
	case int:
		return s, nil
	default:
		return 0, fmt.Errorf("%s: expected string or int, got %T", key, s)
	}
}

func (p *Params) Bool(key string) (bool, error) {
	v := p.Get(key)
	switch s := v.(type) {
	case string:
		return strconv.ParseBool(s)
	case bool:
		return s, nil
	default:
		return false, fmt.Errorf("%s: expected string or bool, got %T", key, s)
	}
}

func (p *Params) Float64(key string) (float64, error) {
	v := p.Get(key)
	switch s := v.(type) {
	case string:
		return strconv.ParseFloat(s, 64)
	case float64:
		return s, nil
	default:
		return 0, fmt.Errorf("%s: expected string or float64, got %T", key, s)
	}
}

func (p *Params) Date(key string, format string) (time.Time, error) {
	v := p.Get(key)
	switch s := v.(type) {
	case string:
		return time.Parse(format, s)
	case time.Time:
		return s, nil
	default:
		return time.Time{}, fmt.Errorf("%s: expected string or time.Time, got %T", key, s)
	}
}

func (p *Params) ParseQuery(query string) error {
	u, err := url.ParseQuery(query)
	if err != nil {
		return err
	}
	for key, values := range u {
		p.data[key] = values
	}
	return nil
}

func (p *Params) AggregateInts(key string, op func(int, int) int) (int, error) {
	vals, ok := p.data[key].([]string)
	if !ok {
		return 0, fmt.Errorf("%s: expected []string, got %T", key, p.data[key])
	}
	var sum int
	for _, val := range vals {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		sum = op(sum, i)
	}
	return sum, nil
}

func (p *Params) SerializeJSON() ([]byte, error) {
	return json.Marshal(p.data)
}

func (p *Params) DeserializeJSON(b []byte) error {
	return json.Unmarshal(b, &p.data)
}

func main() {
	params := NewParams()
	query := "name=Alice&age=25&is_active=true&tags=sport,music&gender=female&pagination=2&sort=asc&dates=2023-10-01,2023-10-02"
	if err := params.ParseQuery(query); err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	name := params.String("name")
	age, _ := params.Int("age")
	isActive, _ := params.Bool("is_active")
	dates, _ := params.Date("dates", "2006-01-02")

	fmt.Println("\nParsed Data:")
	fmt.Println("Name:", name)
	fmt.Println("Age:", age)
	fmt.Println("IsActive:", isActive)
	fmt.Println("Dates:", dates)

	fmt.Println("\nAggregated Data:")
	tagSum, err := params.AggregateInts("tags", func(a, b int) int {
		return a + b
	})
	if err != nil {
		fmt.Println("Error aggregating tags:", err)
	} else {
		fmt.Println("Tag sum:", tagSum)
	}

	fmt.Println("\nPagination and Sorting:")
	pagination, _ := params.Int("pagination")
	sort, _ := params.String("sort")
	fmt.Println("Pagination:", pagination)
	fmt.Println("Sorting:", sort)

	fmt.Println("\nSerialization and Deserialization:")
	bytes, err := params.SerializeJSON()
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
	} else {
		fmt.Println("Serialized JSON:", string(bytes))
	}

	err = params.DeserializeJSON(bytes)
	if err != nil {
		fmt.Println("Error deserializing from JSON:", err)
	} else {
		fmt.Println("Deserialized Data:", params.data)
	}
}