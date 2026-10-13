# Project 9: Full-Stack Task Manager (Gin + React)

## Description
This project is the exercise for Lesson 10: a complete task manager where a React frontend talks to a Gin REST API that stores tasks in SQLite.

Features: list, add, edit, mark as done, and delete tasks; filter by All / Active / Done; client-side and server-side validation; loading and error states.

## Project Structure
```
project9-fullstack/
├── backend/                 Go + Gin + SQLite
│   ├── main.go              router, CORS middleware, serving the React build
│   ├── handlers.go          HTTP handlers, validation, JSON errors
│   └── db.go                Task model and CRUD (database/sql)
└── frontend/                React + Vite
    ├── vite.config.js       dev proxy: /api -> http://localhost:8080
    └── src/
        ├── api.js           fetch wrapper for the REST API
        ├── App.jsx          state, loading/error handling, filters
        └── components/
            ├── TaskForm.jsx form + validation
            └── TaskItem.jsx single task: toggle, edit, delete
```

## API
| Method | Path                | Body                                  | Response          |
|--------|---------------------|---------------------------------------|-------------------|
| GET    | `/api/tasks?status=all\|active\|done` | —                   | `200` list        |
| GET    | `/api/tasks/:id`    | —                                     | `200` / `404`     |
| POST   | `/api/tasks`        | `{"title", "description"}`            | `201` / `400`     |
| PUT    | `/api/tasks/:id`    | `{"title", "description", "done"}`    | `200` / `400` / `404` |
| DELETE | `/api/tasks/:id`    | —                                     | `204` / `404`     |
| GET    | `/api/health`       | —                                     | `200`             |

Errors are returned as `{"error": "message"}`. Validation: `title` is required (max 100 chars), `description` max 500 chars.

## How It Works
1. **Development**: two servers run side by side. Vite (port 5173) serves React with hot reload and **proxies** every `/api` request to Gin (port 8080). The browser only sees one origin, so no CORS is needed.
2. **CORS**: `corsMiddleware` in `main.go` still allows `http://localhost:5173`, so you can also call the API directly (e.g. `fetch('http://localhost:8080/api/tasks')`) and see how preflight `OPTIONS` requests work.
3. **Production**: `npm run build` creates `frontend/dist`. Gin serves `/assets/*` as static files and returns `index.html` for every other non-API route — a single binary serves both the API and the UI.
4. The SQLite driver is `modernc.org/sqlite` (pure Go), so no C compiler / CGO is required — unlike `mattn/go-sqlite3` from Lesson 7.

## How to Run

### Development mode (two terminals)
```bash
cd project9-fullstack/backend
go run .
```
```bash
cd project9-fullstack/frontend
npm install
npm run dev
```
Open http://localhost:5173

### Production mode (one server)
```bash
cd project9-fullstack/frontend
npm install
npm run build
```
```bash
cd project9-fullstack/backend
go build -o taskmanager.exe .
./taskmanager.exe
```
Open http://localhost:8080

### Configuration (environment variables)
| Variable        | Default                  |
|-----------------|--------------------------|
| `PORT`          | `8080`                   |
| `DB_PATH`       | `./tasks.db`             |
| `CORS_ORIGIN`   | `http://localhost:5173`  |
| `FRONTEND_DIST` | `../frontend/dist`       |

### Testing with curl
```bash
curl -X POST http://localhost:8080/api/tasks -H "Content-Type: application/json" -d '{"title":"Learn Gin","description":"Lesson 10"}'
curl http://localhost:8080/api/tasks
curl -X PUT http://localhost:8080/api/tasks/1 -H "Content-Type: application/json" -d '{"title":"Learn Gin","description":"Lesson 10","done":true}'
curl "http://localhost:8080/api/tasks?status=done"
curl -X DELETE http://localhost:8080/api/tasks/1
```

## Notes
- `create-react-app` (mentioned in the lesson plan) is deprecated; Vite is the current standard. Its `server.proxy` replaces the `"proxy"` field from CRA's `package.json`.
- For deployment, set `GIN_MODE=release` and ship the Go binary together with the `frontend/dist` folder (or embed it with `go:embed`).
