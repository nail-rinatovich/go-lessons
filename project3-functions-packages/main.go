package main

import (
	"fmt"
	"project3-functions-packages/math"
)

func main() {
	fmt.Println("Mathematical Operations Demo")
	fmt.Println("============================")

	// Addition
	result := math.Add(10, 5)
	fmt.Printf("10 + 5 = %d\n", result)

	// Subtraction
	result = math.Subtract(10, 5)
	fmt.Printf("10 - 5 = %d\n", result)

	// Multiplication
	result = math.Multiply(10, 5)
	fmt.Printf("10 * 5 = %d\n", result)

	// Division
	result, err := math.Divide(10, 5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 5 = %d\n", result)
	}

	// Division by zero
	result, err = math.Divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 0 = %d\n", result)
	}
}
