package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func sanitizeQueryParam(param string) string {
	// Remove any HTML or JavaScript tags
	param = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(param, "")

	// Remove any single quotes
	param = strings.Replace(param, "'", "\\'", -1)

	// Remove any double quotes
	param = strings.Replace(param, "\"", "\\\"", -1)

	// Remove any backslashes (except for those used to escape quotes)
	param = strings.Replace(param, `\`, "\\", -1)

	return param
}

func main() {
	// Example query string
	queryString := "user=admin&password=secret&action=<script>alert('XSS');</script>"

	// Parse the query string
	parsedQuery, err := url.ParseQuery(queryString)
	if err != nil {
		fmt.Println("Error parsing query string:", err)
		return
	}

	// Sanitize each parameter
	sanitizedParams := make(url.Values)
	for key, value := range parsedQuery {
		sanitizedValues := make([]string, len(value))
		for i, v := range value {
			sanitizedValues[i] = sanitizeQueryParam(v)
		}
		sanitizedParams[key] = sanitizedValues
	}

	// Convert back to a query string
	sanitizedQueryString := sanitizedParams.Encode()

	fmt.Println("Original query string:", queryString)
	fmt.Println("Sanitized query string:", sanitizedQueryString)
}