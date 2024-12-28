package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

// FirewallRule defines the structure of each firewall rule.
type FirewallRule struct {
	Pattern  string `json:"pattern"`
	RuleType string `json:"ruleType"`
	Enabled  bool   `json:"enabled"`
}

// Global variable to hold the current firewall rules
var rules []FirewallRule

// Main function to start the HTTP server and handlers
func main() {
	// Set up logging
	log.SetOutput(os.Stdout)

	// Initial rules setup
	rules = []FirewallRule{
		{Pattern: "DROP TABLE", RuleType: "SQL Injection", Enabled: true},
		{Pattern: "<script>", RuleType: "XSS", Enabled: true},
	}

	// Set up HTTP server with routes
	http.HandleFunc("/update-rules", updateRulesHandler)
	http.HandleFunc("/check-request", checkRequestHandler)

	// Start the server
	log.Println("Firewall service is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// updateRulesHandler updates the current firewall rules dynamically.
func updateRulesHandler(w http.ResponseWriter, r *http.Request) {
	var updatedRules []FirewallRule
	err := json.NewDecoder(r.Body).Decode(&updatedRules)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the rule update to the terminal
	log.Printf("Received rule update: %+v\n", updatedRules)

	// Update the global rules
	rules = updatedRules

	// Log the audit event of rule update
	logAuditEvent("Firewall rules updated")

	// Respond with the updated rules
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedRules)
}

// checkRequestHandler checks if a request matches any of the firewall rules and blocks it if needed.
func checkRequestHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("request")

	// Log the incoming request
	log.Printf("Received request to check: %s\n", request)

	blocked := false

	// Check each rule for SQL Injection, XSS, etc.
	for _, rule := range rules {
		if strings.Contains(request, rule.Pattern) && rule.Enabled {
			log.Printf("Request blocked due to rule: %s\n", rule.RuleType)
			blocked = true
			break
		}
	}

	if blocked {
		log.Println("Request blocked by firewall")
		http.Error(w, "Request blocked by firewall", http.StatusForbidden)
	} else {
		log.Println("Request allowed")
		w.Write([]byte("Request allowed"))
	}
}

// logAuditEvent logs any audit event to the terminal and can be extended to log to a file.
func logAuditEvent(event string) {
	log.Printf("Audit event: %s\n", event)
	// Here, you could write the event to an audit file if needed
	// For simplicity, we just log it to the terminal
}