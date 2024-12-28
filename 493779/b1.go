
package main  
import (  
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strings"

    "github.com/gorilla/mux"
    "github.com/microcosm-cc/bluemonday"
)

// FirewallMiddleware is a middleware function that filters and secures query parameters
func FirewallMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get query parameters
        queryParams := r.URL.Query()

        // Apply filtering to each query parameter
        for key, values := range queryParams {
            for i, value := range values {
                filteredValue := filterQueryParameter(value)
                queryParams.Set(key, filteredValue)
                if filteredValue != value {
                    logSuspiciousActivity("XSS or SQL Injection attempt detected:", r)
                }
            }
        }

        // Update the request with the filtered query parameters
        r.URL.RawQuery = queryParams.Encode()

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}

// filterQueryParameter applies filtering to a single query parameter
func filterQueryParameter(param string) string {
    // Filter out suspicious characters for SQL injection and XSS
    filteredParam := regexp.MustCompile(`[^a-zA-Z0-9_ %.-]+`).ReplaceAllString(param, "")
    filteredParam = strings.TrimSpace(filteredParam)

    // Apply XSS protection using bluemonday
    p := bluemonday.NewPolicy()
    filteredParam = p.Sanitize(filteredParam)

    return filteredParam
}

// logSuspiciousActivity logs suspicious activity to the console
func logSuspiciousActivity(message string, r *http.Request) {
    log.Printf("%s %s from %s\n", message, r.URL.Path, r.RemoteAddr)
}

func main() {
    router := mux.NewRouter()

    // Sample route that handles GET requests to /users
    router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        userID := r.URL.Query().Get("id")
        fmt.Fprintf(w, "User ID: %s\n", userID)
    }).Methods("GET")

    // Apply the firewall middleware to all routes
    http.ListenAndServe(":8080", FirewallMiddleware(router))
}
  