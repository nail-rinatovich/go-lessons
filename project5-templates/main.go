package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Post represents a blog post for demonstration
type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
	Date    string
}

// samplePosts contains some sample data
var samplePosts = []Post{
	{ID: 1, Title: "Getting Started with Go", Content: "Go is a statically typed, compiled programming language designed by Google.", Author: "Alice", Date: "2024-01-15"},
	{ID: 2, Title: "Web Development with Go", Content: "Go's standard library makes it easy to build web servers and APIs.", Author: "Bob", Date: "2024-02-20"},
	{ID: 3, Title: "Templates in Go", Content: "Go's html/template package provides a secure way to generate HTML dynamically.", Author: "Charlie", Date: "2024-03-10"},
}

// funcs are custom functions available inside templates
var funcs = template.FuncMap{
	"upper": strings.ToUpper,
	// excerpt shortens text to n characters
	"excerpt": func(s string, n int) string {
		r := []rune(s) // count characters, not bytes
		if len(r) <= n {
			return s
		}
		return string(r[:n]) + "…"
	},
}

// Templates are parsed once at startup. template.Must stops the program
// if a template file is missing or has a syntax error.
var (
	indexTmpl = template.Must(template.New("index.html").Funcs(funcs).ParseFiles("templates/index.html"))
	postTmpl  = template.Must(template.New("post.html").Funcs(funcs).ParseFiles("templates/post.html"))
)

// findPost returns the post with the given ID
func findPost(id int) (Post, bool) {
	for _, p := range samplePosts {
		if p.ID == id {
			return p, true
		}
	}
	return Post{}, false
}

// indexHandler handles the root path, serving the main page with all posts.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// "/" matches every path, so reject anything that is not the home page
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Execute the template with the posts data
	if err := indexTmpl.Execute(w, samplePosts); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

// postHandler handles requests for individual posts by ID: /post?id=1
func postHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, ok := findPost(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := postTmpl.Execute(w, post); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

func main() {
	// Register handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)

	fmt.Println("Template server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
