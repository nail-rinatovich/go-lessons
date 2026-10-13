# Project 5: Templates

## Description
A small blog rendered with the `html/template` package (Lesson 5).

- `GET /` — list of posts (`templates/index.html`)
- `GET /post?id=1` — a single post (`templates/post.html`)

## How It Works
1. Templates are parsed once at startup with `template.Must`, so a missing file stops the program immediately instead of failing on every request.
2. `index.html` uses `range` to loop over posts and `if` to show a message when the list is empty.
3. Custom template functions are registered with `Funcs`: `upper` (`{{.Author | upper}}` — a pipeline) and `excerpt` (`{{excerpt .Content 50}}`).
4. `postHandler` reads the `id` query parameter: `400` for a non-numeric id, `404` for an unknown post.
5. `html/template` auto-escapes all values, which prevents XSS.

## How to Run
The program reads `templates/` relative to the current directory, so run it from the project folder:
```bash
cd path/to/project5-templates
go run .
```
Open `http://localhost:8080`.
