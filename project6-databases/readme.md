# Project 6: Databases (SQLite)

## Description
A task list stored in SQLite using `database/sql` (Lesson 7). The data is kept in `tasks.db` and survives restarts.

- `GET /tasks` — list all tasks (JSON)
- `POST /tasks` — create a task from form fields `title`, `description`
- `PUT /tasks/{id}` — set `done=true|false`
- `DELETE /tasks/{id}` — delete a task

## How It Works
1. The driver is `modernc.org/sqlite`, a pure-Go SQLite. The popular `github.com/mattn/go-sqlite3` needs CGO and a C compiler (gcc); without them it compiles but fails at startup.
2. `initDB` creates the table and seeds two tasks using a prepared statement (`db.Prepare`).
3. All queries use `?` placeholders — values are never concatenated into SQL, which prevents SQL injection.
4. `db.SetMaxOpenConns(1)` configures the connection pool (SQLite allows one writer at a time).
5. Responses are encoded with `encoding/json`, so titles with quotes still produce valid JSON.
6. Updating or deleting a missing task returns `404`.

## How to Run
```bash
cd path/to/project6-databases
go run .
```
```bash
curl http://localhost:8080/tasks
curl -X POST http://localhost:8080/tasks -d "title=Read a book" -d "description=Chapter 1"
curl -X PUT http://localhost:8080/tasks/1 -d "done=true"
curl -X DELETE http://localhost:8080/tasks/1
```
