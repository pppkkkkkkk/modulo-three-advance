package main

import (
	"fmt"
	"modulo_three_advanced/mod3" 
)

func main() {
	// A basic demonstration of the API
	binaryInput := "1101" // Represents 13
	
	remainder := mod3.ModThree(binaryInput)
	
	fmt.Printf("Input binary: %s\n", binaryInput)
	fmt.Printf("Remainder mod3: %d Expected: 1 \n", remainder) // Expected: 1
}