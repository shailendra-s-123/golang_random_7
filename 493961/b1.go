
package main  
import (  
    "fmt"
    "net/http"
    "sync"
)

type Event string
type Handler func(w http.ResponseWriter, r *http.Request)

var (  
    eventHandlers = make(map[Event]Handler)
    wg            sync.WaitGroup
)

func registerHandler(event Event, handler Handler) {  
    eventHandlers[event] = handler
}

func dispatchEvent(event Event, w http.ResponseWriter, r *http.Request) {  
    handler, ok := eventHandlers[event]
    if ok {  
        handler(w, r)
    } else {  
        http.Error(w, "Not found", http.StatusNotFound)
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {  
    fmt.Fprintf(w, "Hello, world!")
}

func main() {  
    registerHandler("/hello", helloHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {  
        dispatchEvent(Event(r.URL.Path), w, r)
    })

    wg.Add(1)
    go func() {  
        fmt.Println("Server is running on port 8080")
        err := http.ListenAndServe(":8080", nil)
        if err != nil {  
            fmt.Println("Error starting server:", err)
        }
        wg.Done()
    }()

    wg.Wait()
}
                    