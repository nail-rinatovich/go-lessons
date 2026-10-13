# Project 4: Web Server with net/http

## Description
A web server built with Go's standard `net/http` package (Lesson 4). It shows the current server time and greets the user by name submitted through a form.

- `GET /` — page with the current time and a form
- `POST /greet` — reads the `name` form field and responds with a greeting
- `GET /static/...` — static files from the `static/` directory

## How It Works
1. `helloHandler` renders an HTML page with the current time. Any path other than `/` returns `404`.
2. `greetHandler` accepts only `POST` (`405` otherwise), parses the form, and returns `400` if the name is empty.
3. The name is escaped with `html.EscapeString` before it is put into HTML — this prevents XSS (try submitting `<script>alert(1)</script>`).
4. `http.FileServer` + `http.StripPrefix` serve `static/style.css`.

## How to Run
```bash
cd path/to/project4-web-server
go run .
```
Open `http://localhost:8080`.

```bash
curl -X POST http://localhost:8080/greet -d "name=Alice"
```
