package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler responds with a simple greeting.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! This is a simple Go web server.")
}

func main() {
	// Register the handler for the root path.
	http.HandleFunc("/", helloHandler)

	// Start the server on port 8080.
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}