# Project 7: REST API with Gin

## Description
A REST API for managing books, built with Gin and SQLite (Lesson 8).

| Method | Path         | Auth    | Response              |
|--------|--------------|---------|-----------------------|
| GET    | `/books`     | —       | `200` page of books   |
| GET    | `/books/:id` | —       | `200` / `404`         |
| POST   | `/books`     | API key | `201` + `Location` / `400` / `401` |
| PUT    | `/books/:id` | API key | `200` / `400` / `401` / `404` |
| DELETE | `/books/:id` | API key | `204` / `401` / `404` |

Errors are returned as `{"error": "message"}`.

## Features
- **Validation**: `title`, `author`, `year` are required; `year` must be between 1450 and 2100 (`binding` tags).
- **Pagination**: `?page=2&limit=10` (limit 1–100). The response includes `page`, `limit`, `total`.
- **Filtering**: `?author=Jon Bodner` (exact), `?q=go` (title contains).
- **Sorting**: `?sort=title|author|year|id&order=asc|desc`. Column names cannot be SQL placeholders, so they are checked against a whitelist.
- **Authentication**: write requests require the `X-API-Key` header. The key comes from the `API_KEY` environment variable (default `secret-key`).

## How to Run
```bash
cd path/to/project7-rest-api
go run .
```
```bash
curl "http://localhost:8080/books?sort=year&order=desc&limit=2"
curl -X POST http://localhost:8080/books -H "Content-Type: application/json" -H "X-API-Key: secret-key" -d '{"title":"Let'"'"'s Go","author":"Alex Edwards","year":2023}'
curl -X PUT http://localhost:8080/books/1 -H "Content-Type: application/json" -H "X-API-Key: secret-key" -d '{"title":"The Go Programming Language","author":"Alan Donovan","year":2015}'
curl -X DELETE http://localhost:8080/books/1 -H "X-API-Key: secret-key"
```
