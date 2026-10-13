package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// bookInput is the request body for POST and PUT
type bookInput struct {
	Title  string `json:"title" binding:"required,max=200"`
	Author string `json:"author" binding:"required,max=100"`
	Year   int    `json:"year" binding:"required,min=1450,max=2100"`
	ISBN   string `json:"isbn" binding:"max=20"`
}

// sortColumns is the whitelist of columns allowed in ?sort=
var sortColumns = map[string]bool{"id": true, "title": true, "author": true, "year": true}

// apiKeyAuth requires the X-API-Key header to match the configured key
func apiKeyAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid API key"})
			return
		}
		c.Next()
	}
}

type BookHandler struct {
	store *BookStore
}

// GET /books?page=1&limit=10&author=...&q=...&sort=title&order=asc
func (h *BookHandler) list(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err1 != nil || err2 != nil || page < 1 || limit < 1 || limit > 100 {
		jsonError(c, http.StatusBadRequest, "page must be >= 1, limit must be between 1 and 100")
		return
	}

	sort := c.DefaultQuery("sort", "id")
	if !sortColumns[sort] {
		jsonError(c, http.StatusBadRequest, "sort must be one of: id, title, author, year")
		return
	}
	order := strings.ToLower(c.DefaultQuery("order", "asc"))
	if order != "asc" && order != "desc" {
		jsonError(c, http.StatusBadRequest, "order must be asc or desc")
		return
	}

	books, total, err := h.store.List(ListParams{
		Page: page, Limit: limit,
		Author: c.Query("author"), Search: c.Query("q"),
		Sort: sort, Order: order,
	})
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GET /books/:id
func (h *BookHandler) get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	book, err := h.store.Get(id)
	if err != nil {
		storeError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

// POST /books
func (h *BookHandler) create(c *gin.Context) {
	var in bookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	book, err := h.store.Create(Book{Title: in.Title, Author: in.Author, Year: in.Year, ISBN: in.ISBN})
	if err != nil {
		serverError(c, err)
		return
	}
	c.Header("Location", "/books/"+strconv.FormatInt(book.ID, 10))
	c.JSON(http.StatusCreated, book)
}

// PUT /books/:id
func (h *BookHandler) update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in bookInput
	if err := c.ShouldBindJSON(&in); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}
	book, err := h.store.Update(id, Book{Title: in.Title, Author: in.Author, Year: in.Year, ISBN: in.ISBN})
	if err != nil {
		storeError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

// DELETE /books/:id
func (h *BookHandler) delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.store.Delete(id); err != nil {
		storeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		jsonError(c, http.StatusBadRequest, "invalid book id")
		return 0, false
	}
	return id, true
}

func storeError(c *gin.Context, err error) {
	if errors.Is(err, ErrNotFound) {
		jsonError(c, http.StatusNotFound, err.Error())
		return
	}
	serverError(c, err)
}

func serverError(c *gin.Context, err error) {
	log.Println("internal error:", err)
	jsonError(c, http.StatusInternalServerError, "internal server error")
}

func jsonError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func main() {
	db, err := openDB("./books.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "secret-key"
		log.Println("API_KEY not set, using default key: secret-key")
	}

	h := &BookHandler{store: &BookStore{db: db}}
	r := gin.Default()

	// Reading is public
	r.GET("/books", h.list)
	r.GET("/books/:id", h.get)

	// Writing requires an API key
	protected := r.Group("/books", apiKeyAuth(apiKey))
	protected.POST("", h.create)
	protected.PUT("/:id", h.update)
	protected.DELETE("/:id", h.delete)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
