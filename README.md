# Go Web Programming Learning Plan: 10 Lessons

This document outlines a 10-lesson plan for learning web programming with Go, progressing from the basics to building a full-stack application using Gin (Go backend) and React (frontend).

---

## Lesson 1: Introduction to Go and Setting Up the Environment
- **Goal**: Install Go, configure your workspace, and write your first Go program.
- **Topics**:
  - Downloading and installing Go from golang.org
  - Setting up GOPATH and GOBIN (or using Go modules)
  - Writing, compiling, and running a "Hello, World!" program
  - Using `go run` and `go build`
  - Basic Go syntax: package declaration, imports, and the `main` function
- **Exercise**: Create a program that prints your name and favorite programming language.

## Lesson 2: Go Basics: Syntax, Variables, and Control Structures
- **Goal**: Master Go's fundamental building blocks.
- **Topics**:
  - Variables and constants (`var`, `const`, short declaration `:=`)
  - Basic data types: integers, floats, booleans, strings
  - Control structures: `if/else`, `for` loops (including range), `switch`
  - Type conversion and type inference
  - Working with strings (concatenation, formatting with `fmt`)
- **Exercise**: Write a program that calculates the area of a rectangle given user input.

## Lesson 3: Functions and Packages in Go
- **Goal**: Understand code organization and reusability in Go.
- **Topics**:
  - Defining functions: parameters, return values, named returns
  - Variadic functions and defer statements
  - Creating and using custom packages
  - Introduction to Go modules for dependency management
  - Package visibility (exported vs. unexported identifiers)
  - The `init` function
- **Exercise**: Create a package for mathematical operations (add, subtract, multiply, divide) and use it in a main program.

## Lesson 4: Building Web Servers with Go's net/http Package
- **Goal**: Create a basic HTTP server using Go's standard library.
- **Topics**:
  - The `net/http` package: `http.HandleFunc`, `http.ListenAndServe`
  - Handling different HTTP methods (GET, POST, etc.)
  - Parsing URL query parameters and form data
  - Setting response headers and status codes
  - Serving static files (HTML, CSS, images)
  - Basic error handling in HTTP handlers
- **Exercise**: Build a simple web server that displays the current time and accepts a name via a form to greet the user.

## Lesson 5: Working with Templates in Go
- **Goal**: Generate dynamic HTML content using Go's templating engine.
- **Topics**:
  - The `html/template` package (security-conscious templating)
  - Template syntax: actions (`{{ }}`), pipelines, variables
  - Control structures in templates: `if`, `with`, `range`
  - Template functions and custom functions
  - Parsing templates from files and executing them
  - Preventing XSS with auto-escaping
- **Exercise**: Create a blog-like page that displays a list of posts from a slice, using a template.

## Lesson 6: Introduction to Gin Web Framework
- **Goal**: Learn to build web applications faster with the Gin framework.
- **Topics**:
  - Installing Gin: `go get -u github.com/gin-gonic/gin`
  - Setting up a basic Gin router
  - Handling GET, POST, PUT, DELETE requests
  - Working with JSON: binding request bodies and rendering JSON responses
  - Middleware: logging, authentication, CORS
  - Route groups and parameter handling
  - Serving static files with Gin
- **Exercise**: Rebuild the web server from Lesson 4 using Gin, adding JSON endpoints for data.

## Lesson 7: Connecting to Databases with Go
- **Goal**: Integrate persistent storage into your Go web applications.
- **Topics**:
  - Database/sql package overview
  - Working with a specific driver (e.g., `github.com/lib/pq` for PostgreSQL or `github.com/mattn/go-sqlite3` for SQLite)
  - Connection pooling and managing database connections
  - Executing queries: `Query`, `QueryRow`, `Exec`
  - Preventing SQL injection with prepared statements
  - Introduction to ORMs (brief look at GORM)
  - Performing CRUD operations
- **Exercise**: Create a simple task list application that stores tasks in a SQLite database.

## Lesson 8: Building RESTful APIs with Gin
- **Goal**: Design and implement a RESTful API using Gin and a database.
- **Topics**:
  - REST principles: resources, HTTP methods, status codes
  - Structuring API endpoints (e.g., `/users`, `/users/:id`)
  - Request validation and error handling
  - Returning appropriate HTTP status codes and JSON error responses
  - Implementing pagination, filtering, and sorting
  - Basic authentication (API keys or JWT introduction)
  - Testing API endpoints with tools like curl or Postman
- **Exercise**: Build a REST API for managing a collection of books (CRUD operations) using Gin and SQLite.

## Lesson 9: Introduction to React and Setting Up the Frontend
- **Goal**: Learn the basics of React for building user interfaces.
- **Topics**:
  - Installing Node.js and npm
  - Creating a new React app with `create-react-app`
  - Understanding JSX syntax
  - Components: functional vs. class components (focus on functional with hooks)
  - State management with `useState`
  - Handling events and user input
  - Conditional rendering and lists
  - Basic styling with CSS or CSS modules
  - Fetching data from an API using `fetch` or `axios`
- **Exercise**: Create a simple React component that displays a list of items fetched from a mock API.

## Lesson 10: Full-Stack Application: Gin Backend with React Frontend
- **Goal**: Connect your React frontend to your Gin backend to create a full-stack web application.
- **Topics**:
  - Setting up CORS in Gin to allow requests from the React development server
  - Proxying API requests during development (using `package.json` proxy or middleware)
  - Building and consuming REST API endpoints from React
  - Handling asynchronous data loading and loading states
  - Form handling and validation in React
  - Optimizing builds for production
  - Deployment considerations (e.g., serving React build files from Gin)
  - Full-stack project structure and organization
- **Exercise**: Build a complete task manager application where the React frontend interacts with the Gin backend API to display, add, update, and delete tasks stored in a database.

---

## Next Steps
After completing these 10 lessons, you will have a solid foundation in Go web development and be able to build modern full-stack applications. Consider exploring advanced topics such as:
- WebSocket communication with Gorilla WebSocket
- Authentication and authorization (OAuth2, sessions)
- Testing Go applications (unit tests, integration tests)
- Dockerizing your Go and React applications
- Cloud deployment (AWS, Google Cloud, Heroku)
- Advanced React patterns (context, reducers, routing with React Router)

Happy coding!