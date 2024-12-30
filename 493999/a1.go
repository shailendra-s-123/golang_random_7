package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// QueryParameterManager is a struct to manage URL query parameters
type QueryParameterManager struct {
	params map[string][]interface{}
}

// NewQueryParameterManager returns a new instance of QueryParameterManager
func NewQueryParameterManager() *QueryParameterManager {
	return &QueryParameterManager{params: make(map[string][]interface{})}
}

// ParseQueryString parses a URL query string and adds it to the manager
func (m *QueryParameterManager) ParseQueryString(query string) {
	values, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}

	for key, vals := range values {
		for _, val := range vals {
			m.AddParameter(key, val)
		}
	}
}

// AddParameter adds a parameter to the manager
func (m *QueryParameterManager) AddParameter(key string, value interface{}) {
	m.params[key] = append(m.params[key], value)
}

// GetParameter returns all values for a given key, with type conversion
func (m *QueryParameterManager) GetParameter(key string) ([]interface{}, error) {
	params := m.params[key]
	if params == nil {
		return nil, fmt.Errorf("parameter %q not found", key)
	}
	return params, nil
}

// GetStringParameter returns the first string value for a given key
func (m *QueryParameterManager) GetStringParameter(key string) (string, error) {
	params, err := m.GetParameter(key)
	if err != nil {
		return "", err
	}
	if len(params) == 0 {
		return "", fmt.Errorf("no string value found for parameter %q", key)
	}
	return params[0].(string), nil
}

// GetIntParameter returns the first int value for a given key
func (m *QueryParameterManager) GetIntParameter(key string) (int, error) {
	params, err := m.GetParameter(key)
	if err != nil {
		return 0, err
	}
	if len(params) == 0 {
		return 0, fmt.Errorf("no int value found for parameter %q", key)
	}
	val, ok := params[0].(int)
	if !ok {
		return 0, fmt.Errorf("value for parameter %q is not an int", key)
	}
	return val, nil
}

// GetBoolParameter returns the first bool value for a given key
func (m *QueryParameterManager) GetBoolParameter(key string) (bool, error) {
	params, err := m.GetParameter(key)
	if err != nil {
		return false, err
	}
	if len(params) == 0 {
		return false, fmt.Errorf("no bool value found for parameter %q", key)
	}
	val, ok := params[0].(bool)
	if !ok {
		return false, fmt.Errorf("value for parameter %q is not a bool", key)
	}
	return val, nil
}

func main() {
	query := "start=2023-01-01&end=2023-06-30&useApi=true&limit=100&start=2023-07-01&limit=50"

	m := NewQueryParameterManager()
	m.ParseQueryString(query)

	// Get specific parameters by type
	startDateStr, err := m.GetStringParameter("start")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("First start date:", startDateStr)

	endDateStr, err := m.GetStringParameter("end")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("End date:", endDateStr)

	useApi, err := m.GetBoolParameter("useApi")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Use API:", useApi)

	limit, err := m.GetIntParameter("limit")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Limit:", limit)

	// Get all values for a parameter
	allStartDates, err := m.GetParameter("start")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("All start dates:", allStartDates)

	// Example of filtering based on a parameter (limit)
	limits := []int{}
	for _, param := range m.params["limit"] {
		limit, ok := param.(int)
		if ok {
			limits = append(limits, limit)
		}
	}

	maxLimit := 0
	for _, limit := range limits {
		if limit > maxLimit {
			maxLimit = limit
		}
	}

	fmt.Println("Maximum limit found:", maxLimit)
}