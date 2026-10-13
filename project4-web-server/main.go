package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

// helloHandler responds with a simple greeting and the current time.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// "/" matches every path, so reject anything that is not the home page
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Get current time
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// HTML form to submit a name
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Go Web Server</title>
			<link rel="stylesheet" href="/static/style.css">
		</head>
		<body>
			<h1>Welcome to the Go Web Server!</h1>
			<p>Current server time: %s</p>
			<form action="/greet" method="post">
				<label for="name">Enter your name:</label>
				<input type="text" id="name" name="name" required>
				<button type="submit">Greet Me</button>
			</form>
		</body>
		</html>
	`, currentTime)

	fmt.Fprint(w, html)
}

// greetHandler processes the form submission and responds with a greeting.
func greetHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the name from the form
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Respond with a greeting. Escape the name so user input
	// cannot inject HTML or JavaScript into the page (XSS)
	fmt.Fprintf(w, "<h1>Hello, %s! Nice to meet you.</h1><a href=\"/\">Back</a>", html.EscapeString(name))
}

func main() {
	// Register handlers
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/greet", greetHandler)
	// Serve static files from the "static" directory.
	// Strip the "/static/" prefix before looking for the file
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}