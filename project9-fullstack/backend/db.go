package main

import (
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// ErrNotFound is returned when a task with the given ID does not exist
var ErrNotFound = errors.New("task not found")

// Task represents a task stored in the database
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdAt"`
}

// openDB opens the SQLite database and creates the schema
func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite allows only one writer at a time
	db.SetMaxOpenConns(1)

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		done BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// TaskStore wraps the database and provides CRUD operations
type TaskStore struct {
	db *sql.DB
}

// List returns tasks filtered by status: "all", "active" or "done"
func (s *TaskStore) List(status string) ([]Task, error) {
	query := "SELECT id, title, description, done, created_at FROM tasks"
	switch status {
	case "active":
		query += " WHERE done = 0"
	case "done":
		query += " WHERE done = 1"
	}
	query += " ORDER BY id DESC"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Start with an empty slice so the JSON is [] instead of null
	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// Get returns a single task by ID
func (s *TaskStore) Get(id int64) (Task, error) {
	var t Task
	err := s.db.QueryRow(
		"SELECT id, title, description, done, created_at FROM tasks WHERE id = ?", id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return t, ErrNotFound
	}
	return t, err
}

// Create inserts a new task and returns it
func (s *TaskStore) Create(title, description string) (Task, error) {
	result, err := s.db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		return Task{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}
	return s.Get(id)
}

// Update replaces the title, description and done status of a task
func (s *TaskStore) Update(id int64, title, description string, done bool) (Task, error) {
	result, err := s.db.Exec(
		"UPDATE tasks SET title = ?, description = ?, done = ? WHERE id = ?",
		title, description, done, id,
	)
	if err != nil {
		return Task{}, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return Task{}, ErrNotFound
	}
	return s.Get(id)
}

// Delete removes a task by ID
func (s *TaskStore) Delete(id int64) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}
