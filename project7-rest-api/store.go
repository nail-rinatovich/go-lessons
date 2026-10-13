package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// ErrNotFound is returned when a book with the given ID does not exist
var ErrNotFound = errors.New("book not found")

// Book is the resource managed by the API
type Book struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	ISBN   string `json:"isbn"`
}

// ListParams holds pagination, filtering and sorting options
type ListParams struct {
	Page   int
	Limit  int
	Author string // exact author filter
	Search string // substring search in title
	Sort   string // column name, validated by the handler
	Order  string // "asc" or "desc"
}

// BookStore wraps the database and provides CRUD operations
type BookStore struct {
	db *sql.DB
}

// openDB opens the SQLite database, creates the schema and seeds sample data
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite allows only one writer at a time
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		year INTEGER NOT NULL,
		isbn TEXT NOT NULL DEFAULT ''
	);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count); err != nil {
		db.Close()
		return nil, err
	}
	if count == 0 {
		seed := []Book{
			{Title: "The Go Programming Language", Author: "Alan Donovan", Year: 2015, ISBN: "978-0134190440"},
			{Title: "Learning Go", Author: "Jon Bodner", Year: 2021, ISBN: "978-1492077213"},
			{Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Year: 2017, ISBN: "978-1491941195"},
		}
		for _, b := range seed {
			if _, err := db.Exec("INSERT INTO books (title, author, year, isbn) VALUES (?, ?, ?, ?)",
				b.Title, b.Author, b.Year, b.ISBN); err != nil {
				db.Close()
				return nil, err
			}
		}
	}
	return db, nil
}

// List returns one page of books and the total number of matching books
func (s *BookStore) List(p ListParams) ([]Book, int, error) {
	// Build the WHERE clause from placeholders only: values are never
	// concatenated into SQL, which prevents SQL injection
	var where []string
	var args []any
	if p.Author != "" {
		where = append(where, "author = ?")
		args = append(args, p.Author)
	}
	if p.Search != "" {
		where = append(where, "title LIKE ?")
		args = append(args, "%"+p.Search+"%")
	}
	whereSQL := ""
	if len(where) > 0 {
		whereSQL = " WHERE " + strings.Join(where, " AND ")
	}

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM books"+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Column names cannot be placeholders, so p.Sort and p.Order must be
	// validated against a whitelist before reaching this point
	query := fmt.Sprintf("SELECT id, title, author, year, isbn FROM books%s ORDER BY %s %s LIMIT ? OFFSET ?",
		whereSQL, p.Sort, p.Order)
	args = append(args, p.Limit, (p.Page-1)*p.Limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.ISBN); err != nil {
			return nil, 0, err
		}
		books = append(books, b)
	}
	return books, total, rows.Err()
}

// Get returns a single book by ID
func (s *BookStore) Get(id int64) (Book, error) {
	var b Book
	err := s.db.QueryRow("SELECT id, title, author, year, isbn FROM books WHERE id = ?", id).
		Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.ISBN)
	if errors.Is(err, sql.ErrNoRows) {
		return b, ErrNotFound
	}
	return b, err
}

// Create inserts a new book and returns it with its ID
func (s *BookStore) Create(b Book) (Book, error) {
	result, err := s.db.Exec("INSERT INTO books (title, author, year, isbn) VALUES (?, ?, ?, ?)",
		b.Title, b.Author, b.Year, b.ISBN)
	if err != nil {
		return Book{}, err
	}
	b.ID, err = result.LastInsertId()
	return b, err
}

// Update replaces all fields of an existing book
func (s *BookStore) Update(id int64, b Book) (Book, error) {
	result, err := s.db.Exec("UPDATE books SET title = ?, author = ?, year = ?, isbn = ? WHERE id = ?",
		b.Title, b.Author, b.Year, b.ISBN, id)
	if err != nil {
		return Book{}, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return Book{}, ErrNotFound
	}
	b.ID = id
	return b, nil
}

// Delete removes a book by ID
func (s *BookStore) Delete(id int64) error {
	result, err := s.db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
