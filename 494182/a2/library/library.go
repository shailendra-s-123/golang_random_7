package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simulate an error
	err := errors.New("sample error")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello, World!\n")
}

func main() {
	http.HandleFunc("/sample", handleRequest)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", errors.Wrap(err, "error starting server"))
	}
}