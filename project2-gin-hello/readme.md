# Project 2: Gin Web Framework

## Description
This project demonstrates building a web server using the Gin framework (Lesson 6). It provides RESTful JSON endpoints for managing blog posts:

- `GET /` — Simple greeting message
- `GET /posts` — Returns a list of all posts
- `POST /posts` — Creates a new post from JSON body

## How It Works
1. Gin router is initialized with default middleware (logging + recovery).
2. Three routes are registered for handling GET and POST requests.
3. The `/posts` endpoint returns a list of posts as JSON.
4. The `POST /posts` endpoint parses the JSON request body, validates it (`title` and `body` are required), assigns a new ID, appends the post to the in-memory slice and returns `201 Created`.
5. A `sync.Mutex` protects the slice, because Gin handles requests concurrently.

## How to Run
1. Navigate to this project's directory:
   ```bash
   cd path/to/project2-gin-hello
   ```
2. Run the program:
   ```bash
   go run main.go
   ```
3. Test the endpoints:
   ```bash
   # Get the greeting
   curl http://localhost:8080/

   # Get all posts
   curl http://localhost:8080/posts

   # Create a new post
   curl -X POST http://localhost:8080/posts \
     -H "Content-Type: application/json" \
     -d '{"title":"New Post","body":"Post content here"}'
   ```

## Expected Output (Terminal)
```
[GIN-debug] GET    /                         --> main.helloHandler
[GIN-debug] GET    /posts                    --> main.postsHandler
[GIN-debug] POST   /posts                    --> main.createPostHandler
Server starting on http://localhost:8080
```

## Notes
- Uses in-memory storage for demonstration (data is lost on restart).
- Gin provides built-in middleware for logging, recovery, and JSON binding.