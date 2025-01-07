package main

import (
	"fmt"
	"net/http"
)

// ErrorHandler is a middleware that handles errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				switch err := r.(type) {
				case error:
					http.Error(w, err.Error(), http.StatusInternalServerError)
				default:
					http.Error(w, fmt.Sprintf("%v", r), http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", ErrorHandler(http.DefaultServeMux))
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Simulate an error
    panic("Something went wrong!")
}