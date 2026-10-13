# Project 1: Simple Go Web Server

## Description
This project demonstrates a basic HTTP server using Go's standard `net/http` package. It listens on port 8080 and responds with a simple greeting message for any request to the root path ("/").

## How It Works
1. The `main` function sets up an HTTP handler for the root path (`/`) using `http.HandleFunc`.
2. The handler function `helloHandler` writes a plain text response: "Hello, World! This is a simple Go web server."
3. The server starts listening on `:8080` with `http.ListenAndServe`.
4. If the server encounters an error (e.g., port already in use), it logs the error and exits.

## How to Run
1. Ensure you have Go installed (version 1.16 or later recommended).
2. Navigate to this project's directory:
   ```bash
   cd path/to/project1-hello-web
   ```
3. Run the program:
   ```bash
   go run main.go
   ```
4. Open your web browser and visit: `http://localhost:8080`
   You should see the greeting message.
5. To stop the server, press `Ctrl+C` in the terminal where it's running.

## Expected Output
- Terminal: "Server starting on http://localhost:8080"
- Browser: "Hello, World! This is a simple Go web server."

## Notes
- This is a minimal example for learning purposes. In a real application, you would likely handle multiple routes, use middleware, and possibly use a framework like Gin for more features.