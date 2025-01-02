package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Define an EventBus type as a map of event names to a slice of callbacks
type EventBus map[string][]func(interface{})

var (
	// Global event bus
	eventBus EventBus
	// Mutex to protect concurrent access to the event bus
	busMutex sync.Mutex
)

// Register an API endpoint as a callback for a specific event
func RegisterCallback(event string, callback func(interface{})) {
	busMutex.Lock()
	defer busMutex.Unlock()

	if eventBus == nil {
		eventBus = make(EventBus)
	}

	eventBus[event] = append(eventBus[event], callback)
}

// Trigger an event and call all the registered callbacks for that event
func TriggerEvent(event string, data interface{}) {
	busMutex.Lock()
	defer busMutex.Unlock()

	if callbacks, ok := eventBus[event]; ok {
		for _, callback := range callbacks {
			go callback(data) // Call the callback in a goroutine to avoid blocking the event loop
		}
	}
}

// Sample API endpoint callback
func sampleAPIEndpointCallback(data interface{}) {
	eventData := data.(map[string]interface{})
	message := eventData["message"].(string)
	fmt.Println("Received event:", message)
}

func main() {
	// Register the callback for the "new_message" event
	RegisterCallback("new_message", sampleAPIEndpointCallback)

	// Create a simple HTTP server to trigger the event
	http.HandleFunc("/trigger-event", func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Query().Get("message")
		eventData := map[string]interface{}{
			"message": message,
		}
		// Trigger the "new_message" event with the event data
		TriggerEvent("new_message", eventData)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Server is running on port 8080. Press Ctrl+C to exit.")
	fmt.Println("To trigger the event, use: curl 'http://localhost:8080/trigger-event?message=Hello%20World!'")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}  