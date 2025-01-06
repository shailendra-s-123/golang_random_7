
package main  
import (  
    "fmt"
    "net/url"
    "strings"
	"strconv"
)

// QueryParams struct encapsulates the query parameters for data insights
type QueryParams struct {  
    Filters       map[string][]string // Filter criteria
    SortField     string               // Field for sorting
    SortOrder     string               // Sort order (asc/desc)
    Page          int                  // Page number
    PageSize      int                  // Page size
    ExtraParams   url.Values           // Extra parameters
}  

// NewQueryParams creates a new QueryParams instance from URL values
func NewQueryParams(values url.Values) *QueryParams {  
    params := &QueryParams{  
        Filters:    make(map[string][]string),  
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
            params.Filters[field] = values  
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

// ApplyToQuery applies the query parameters to a url.Values instance
func (p *QueryParams) ApplyToQuery(query url.Values) {  
    // Add predefined query parameters
    addParam(query, "sortfield", p.SortField)  
    addParam(query, "sortorder", p.SortOrder)  
    addParam(query, "page", strconv.Itoa(p.Page))  
    addParam(query, "pagesize", strconv.Itoa(p.PageSize))  
    
    // Add filter parameters
    for field, values := range p.Filters {  
        for _, value := range values {  
            addParam(query, fmt.Sprintf("filter_%s", field), value)  
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
    rawQuery := "sortfield=name&sortorder=asc&page=2&pagesize=15&filter_country=US&filter_status=active"
    values, _ := url.ParseQuery(rawQuery)  
    params := NewQueryParams(values)  
    
    fmt.Printf("Filters: %v\n", params.Filters)  
    fmt.Printf("SortField: %s\n", params.SortField)  
    fmt.Printf("SortOrder: %s\n", params.SortOrder)  
    fmt.Printf("Page: %d\n", params.Page)  
    fmt.Printf("PageSize: %d\n", params.PageSize)  
    
    // Modify and re-apply the parameters
}