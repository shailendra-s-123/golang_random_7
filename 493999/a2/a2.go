package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"strconv"
)

// QueryParams struct encapsulates the query parameters for data insights
type QueryParams struct {
	Filters         map[string]interface{} // Filter criteria (now supports various types)
	SortField       string                   // Field for sorting
	SortOrder       string                   // Sort order (asc/desc)
	Page            int                      // Page number
	PageSize        int                      // Page size
	ExtraParams     url.Values               // Extra parameters
}

// NewQueryParams creates a new QueryParams instance from URL values
func NewQueryParams(values url.Values) *QueryParams {
	params := &QueryParams{
		Filters:    make(map[string]interface{}),
		ExtraParams: make(url.Values),
	}

	// Parse predefined query parameters
	for key, values := range values {
		switch strings.ToLower(key) {
		case "sortfield":
			params.SortField = values[0]
		case "sortorder":
			params.SortOrder = values[0]
		case "page":
			params.Page = parseInt(values[0], 1)
		case "pagesize":
			params.PageSize = parseInt(values[0], 10)
		default:
			// Handle extra parameters
			params.ExtraParams[key] = values
		}
	}

	// Parse filter parameters (multi-valued)
	for key, values := range values {
		if strings.HasPrefix(key, "filter_") {
			field := strings.TrimPrefix(key, "filter_")

			if strings.HasPrefix(field, "range_") {
				rangeField := strings.TrimPrefix(field, "range_")
				rangeValues := parseRange(values)
				params.Filters[rangeField] = rangeValues
			} else if strings.HasPrefix(field, "date_") {
				dateField := strings.TrimPrefix(field, "date_")
				dateValue := parseDate(values[0])
				params.Filters[dateField] = dateValue
			} else {
				params.Filters[field] = values
			}
		}
	}

	return params
}

func parseInt(s string, defaultValue int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return value
}

func parseRange(values []string) []int {
	var rangeValues []int
	for _, value := range values {
		if value != "" {
			number, _ := strconv.Atoi(value)
			rangeValues = append(rangeValues, number)
		}
	}
	return rangeValues
}

func parseDate(dateString string) *time.Time {
	if dateString != "" {
		dateLayout := "2006-01-02" // You can adjust this based on your date format
		date, err := time.Parse(dateLayout, dateString)
		if err != nil {
			fmt.Printf("Invalid date format: %s\n", dateString)
			return nil
		}
		return &date
	}
	return nil
}

// ApplyToQuery applies the query parameters to a url.Values instance
func (p *QueryParams) ApplyToQuery(query url.Values) {
	// Add predefined query parameters
	addParam(query, "sortfield", p.SortField)
	addParam(query, "sortorder", p.SortOrder)
	addParam(query, "page", strconv.Itoa(p.Page))
	addParam(query, "pagesize", strconv.Itoa(p.PageSize))

	// Add filter parameters
	for field, value := range p.Filters {
		switch v := value.(type) {
		case []int:
			for _, num := range v {
				addParam(query, fmt.Sprintf("filter_range_%s", field), strconv.Itoa(num))
			}
		case *time.Time:
			if v != nil {
				dateLayout := "2006-01-02"
				addParam(query, fmt.Sprintf("filter_date_%s", field), v.Format(dateLayout))
			}
		default:
			addParam(query, fmt.Sprintf("filter_%s", field), v.([]string)[0])
		}
	}

	// Add extra parameters
	for key, values := range p.ExtraParams {
		query[key] = values
	}
}

func addParam(query url.Values, key, value string) {
	if value != "" {
		query.Add(key, value)
	}
}

func main() {
	// Example usage
	rawQuery := "sortfield=name&sortorder=asc&page=2&pagesize=15&filter_country=US&filter_status=active&filter_range_age=20,30&filter_date_start=2023-01-01"
	values, _ := url.ParseQuery(rawQuery)
	params := NewQueryParams(values)

	fmt.Printf("Filters: %v\n", params.Filters)
	fmt.Printf("SortField: %s\n", params.SortField)
	fmt.Printf("SortOrder: %s\n", params.SortOrder)
	fmt.Printf("Page: %d\n", params.Page)
	fmt.Printf("PageSize: %d\n", params.PageSize)
}