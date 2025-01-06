package main

import (
	_"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"strconv"
)

// QueryParameters is a struct to manage URL query parameters.
type QueryParameters struct {
	Filters    []*Filter  `json:"filters,omitempty"`
	Sort       *Sort      `json:"sort,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Filter is a struct to represent a filter condition.
type Filter struct {
	Key    string
	Value  string
	Operator string
}

// Sort is a struct to represent a sorting condition.
type Sort struct {
	Key    string
	Order  string // e.g., "asc" or "desc"
}

// Pagination is a struct to represent pagination details.
type Pagination struct {
	Page     int
	PageSize int
}

// ParseQuery from URL.Values
func (q *QueryParameters) ParseQuery(qr url.Values) error {
	var err error
	q.Filters, err = parseFilters(qr["filters"])
	if err != nil {
		return err
	}
	q.Sort, err = parseSort(qr["sort"])
	if err != nil {
		return err
	}
	q.Pagination, err = parsePagination(qr["pagination"])
	if err != nil {
		return err
	}
	return nil
}

func parseFilters(raw []string) ([]*Filter, error) {
	filters := make([]*Filter, 0)
	for _, rawFilter := range raw {
		parts := strings.Split(rawFilter, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid filter format: %s", rawFilter)
		}
		filter := &Filter{
			Key:    strings.Trim(parts[0], " "),
			Value:  strings.Trim(parts[1], " "),
			Operator: strings.Trim(parts[2], " "),
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func parseSort(raw []string) (*Sort, error) {
	if len(raw) != 1 {
		return nil, fmt.Errorf("invalid sort format: %s", raw)
	}
	parts := strings.Split(raw[0], ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid sort format: %s", raw)
	}
	return &Sort{
		Key:    strings.Trim(parts[0], " "),
		Order:  strings.Trim(parts[1], " "),
	}, nil
}

func parsePagination(raw []string) (*Pagination, error) {
	if len(raw) != 1 {
		return nil, fmt.Errorf("invalid pagination format: %s", raw)
	}
	parts := strings.Split(raw[0], ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid pagination format: %s", raw)
	}
	page, err := strconv.Atoi(strings.Trim(parts[0], " "))
	if err != nil {
		return nil, fmt.Errorf("invalid page number: %s", parts[0])
	}
	pageSize, err := strconv.Atoi(strings.Trim(parts[1], " "))
	if err != nil {
		return nil, fmt.Errorf("invalid page size: %s", parts[1])
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func main() {
	// Example URL with query parameters
	urlStr := "https://api.example.com/data?filters=key1=value1,operator=eq&sort=key2,desc&pagination=1,10"

	// Parse URL and query parameters
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	qr := u.Query()

	// Create a new QueryParameters struct
	qp := &QueryParameters{}
	if err := qp.ParseQuery(qr); err != nil {
		fmt.Println("Error parsing query:", err)
		return
	}

	// Output parsed parameters
	fmt.Println(qp)
}