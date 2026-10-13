# Project 3: Functions and Packages in Go

## Description
This project demonstrates how to create and use custom packages in Go. It includes a math package with basic arithmetic operations (add, subtract, multiply, divide) and a main program that uses this package.

## How It Works
1. The `math` package contains functions for basic arithmetic operations.
2. The `main` program imports the `math` package and demonstrates each operation.
3. The `Divide` function returns an error when attempting to divide by zero.

## How to Run
1. Ensure you have Go installed (version 1.16 or later recommended).
2. Navigate to this project's directory:
   ```bash
   cd path/to/project3-functions-packages
   ```
3. Run the program:
   ```bash
   go run main.go
   ```
4. You should see the output of each mathematical operation.

## Expected Output
```
Mathematical Operations Demo
============================
10 + 5 = 15
10 - 5 = 5
10 * 5 = 50
10 / 5 = 2
Error: division by zero
```

## Notes
- This project demonstrates creating a custom package and importing it in another Go file.
- The math package is located in the `math` directory.
- The main program is in `main.go` at the project root.