package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	// Pure-Go SQLite driver: works without CGO / a C compiler.
	// (github.com/mattn/go-sqlite3 needs gcc and CGO_ENABLED=1.)
	_ "modernc.org/sqlite"
)

// Task represents a task in the task list
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// initDB creates the tasks table and seeds it with sample data
func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		done BOOLEAN NOT NULL DEFAULT 0
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Seed with sample data if table is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// A prepared statement is compiled once and reused for every insert
		stmt, err := db.Prepare("INSERT INTO tasks (title, description, done) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec("Learn Go", "Complete the Go web programming course", false); err != nil {
			return err
		}
		if _, err := stmt.Exec("Build a web server", "Create a REST API with Gin", false); err != nil {
			return err
		}
	}
	return nil
}

// getTasks returns all tasks from the database
func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, description, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty slice instead of nil, so the JSON is [] instead of null
	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// createTask inserts a new task into the database
func createTask(db *sql.DB, title, description string) (int64, error) {
	result, err := db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// updateTask sets the done status of a task.
// Returns sql.ErrNoRows if the task does not exist.
func updateTask(db *sql.DB, id int64, done bool) error {
	result, err := db.Exec("UPDATE tasks SET done = ? WHERE id = ?", done, id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// deleteTask removes a task from the database.
// Returns sql.ErrNoRows if the task does not exist.
func deleteTask(db *sql.DB, id int64) error {
	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// writeJSON encodes v as JSON. encoding/json escapes quotes and special
// characters, so a title like `say "hi"` still produces valid JSON.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("Failed to write JSON:", err)
	}
}

// writeDBError maps database errors to HTTP responses
func writeDBError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	log.Println("Database error:", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Connection pool settings. SQLite allows only one writer at a time.
	db.SetMaxOpenConns(1)

	// Initialize the database
	if err := initDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	fmt.Println("Database initialized. Server starting on http://localhost:8080")

	// Register HTTP handlers
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tasks, err := getTasks(db)
			if err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, tasks)

		case http.MethodPost:
			title := strings.TrimSpace(r.FormValue("title"))
			description := strings.TrimSpace(r.FormValue("description"))
			if title == "" {
				http.Error(w, "Title is required", http.StatusBadRequest)
				return
			}
			id, err := createTask(db, title, description)
			if err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, Task{ID: id, Title: title, Description: description})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/tasks/"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			done := r.FormValue("done") == "true"
			if err := updateTask(db, id, done); err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"id": id, "done": done})

		case http.MethodDelete:
			if err := deleteTask(db, id); err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"deleted": true})

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
