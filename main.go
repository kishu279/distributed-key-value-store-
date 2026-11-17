package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Go Web Server!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	// 1. Create a new request router (ServeMux)
	mux := http.NewServeMux()

	// 2. Register the handler function for the root URL path ("/")
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/home", homeHandler)

	// 3. Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
