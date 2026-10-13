package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Data structure for demonstration
type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// mu protects posts: Gin runs handlers concurrently
var mu sync.Mutex

// In-memory "database" for demonstration purposes
var posts = []Post{
	{ID: 1, Title: "First Post", Body: "This is the first post."},
	{ID: 2, Title: "Second Post", Body: "This is the second post."},
}

// helloHandler responds with a simple greeting.
func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World! This is a Gin web server.",
	})
}

// postsHandler returns a list of all posts as JSON.
func postsHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// createPostHandler creates a new post from the JSON request body.
func createPostHandler(c *gin.Context) {
	var newPost Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Assign a new ID (simple approach for demonstration)
	if len(posts) > 0 {
		newPost.ID = posts[len(posts)-1].ID + 1
	} else {
		newPost.ID = 1
	}

	posts = append(posts, newPost)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created",
		"post":    newPost,
	})
}

func main() {
	// Create a new Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Register handlers
	r.GET("/", helloHandler)
	r.GET("/posts", postsHandler)
	r.POST("/posts", createPostHandler)

	// Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}