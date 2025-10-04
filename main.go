package main

import (
	"fmt"
	"modulo_three_advanced/mod3" 
)

func main() {
	// --- DEMONSTRATION 1: SUCCESS PATH (Valid Input) ---
	binaryInputValid := "1101" // Represents 13 (13 mod 3 = 1)
	
	// 1. Create the calculator service instance.
	calc, err := mod3.NewModThreeCalculator(mod3.GetModThreeConfig())
	if err != nil {
		fmt.Printf("FATAL ERROR: Failed to initialize Modulo Calculator: %v\n", err)
		return
	}
	
	// 2. Call the interface method on the created object.
	remainder, execErr := calc.Calculate(binaryInputValid)
	
	fmt.Printf("--- Test Case 1: Valid Input ---\n")
	fmt.Printf("Input: %q (Decimal 13)\n", binaryInputValid)

	if execErr != nil {
		fmt.Printf("  Result: ERROR Execution \n  Reason: %v\n", execErr)
	} else {
		fmt.Printf("  Result: Success Execution \n  Remainder: %d (Expected: 1)\n", remainder)
	}

	// --- DEMONSTRATION 2: ERROR PATH (Invalid Input) ---
	binaryInputInvalid := "1A01" // Contains invalid character 'A'
	
	remainderInvalid, execErrInvalid := calc.Calculate(binaryInputInvalid)

	fmt.Printf("\n--- Test Case 2: Invalid Input (Error Path) ---\n")
	fmt.Printf("Input: %q (Contains 'A')\n", binaryInputInvalid)

	if execErrInvalid != nil {
		fmt.Printf("  Result: ERROR Execution \n  Remainder: %d (Expected -1 on error) \n  Reason: %v\n", 
			remainderInvalid, execErrInvalid)
	} else {
		fmt.Printf("  Result: Success Execution \n  Remainder: %d\n", remainderInvalid)
	}
}
