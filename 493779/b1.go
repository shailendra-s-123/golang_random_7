package main  
import (  
        "fmt"
        "log"
        "net/http"
        "github.com/gorilla/mux"
        "github.com/microcosm-cc/bluemonday"
        "regexp"
)

// sanitizeInput uses bluemonday to sanitize the given string
func sanitizeInput(input string) string {
        p := bluemonday.NewPolicy()
        p.AllowAttrs().OnElements("a")
        p.RequireParseableURLs(true)
        return p.Sanitize(input)
}

// ValidateEmail uses a regular expression to validate email format
func ValidateEmail(email string) bool {
        const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
        match, _ := regexp.MatchString(emailRegex, email)
        return match
}

func main() {
        r := mux.NewRouter()
        r.HandleFunc("/user", handleUser).Methods("GET")
        log.Fatal(http.ListenAndServe(":8080", r))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
        // Get the URL query parameters
        params := r.URL.Query()
        
        // Validate and sanitize user input parameters
        userID := params.Get("user_id")
        safeUserID := sanitizeInput(userID)
        if safeUserID == "" {
                http.Error(w, "Invalid user ID", http.StatusBadRequest)
                return
        }
        
        email := params.Get("email")
        safeEmail := sanitizeInput(email)
        if safeEmail == "" || !ValidateEmail(safeEmail) {
                http.Error(w, "Invalid email format", http.StatusBadRequest)
                return
        }
        
        // Now you can use the sanitized and validated parameters in your application
        fmt.Fprintf(w, "Hello, User ID: %s, Email: %s", safeUserID, safeEmail)
}