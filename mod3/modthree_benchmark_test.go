package mod3

import (
	"testing"
	"strings"
)

// Benches the ModThreeCalculator.Calculate method for various input lengths.
// Declare a global variable to prevent compiler optimizations
// that might eliminate calculations if results are not used.
var result int

// setupCalculator initializes the ModThreeCalculator once for all benchmarks.
func setupCalculator(b *testing.B) ModuloCalculator {
	cfg := GetModThreeConfig()
	calc, err := NewModThreeCalculator(cfg)
	if err != nil {
		b.Fatalf("Failed to initialize Modulo3 calculator for benchmarks: %v", err)
	}
	return calc
}

// BenchmarkCalculate_Short benchmarks a very short input string.
func BenchmarkCalculate_Short(b *testing.B) {
	calc := setupCalculator(b)
	input := "101" // Binary 5

	b.ResetTimer() // Reset timer after setup
	for i := 0; i < b.N; i++ {
		// Store result in a global to prevent compiler optimization
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_Medium benchmarks a medium-length input string.
func BenchmarkCalculate_Medium(b *testing.B) {
	calc := setupCalculator(b)
	input := strings.Repeat("10", 10) + "1" // 21 bits
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_Long benchmarks a long input string.
func BenchmarkCalculate_Long(b *testing.B) {
	calc := setupCalculator(b)
	input := strings.Repeat("1", 100) // 100 bits
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_VeryLong benchmarks a very long input string.
func BenchmarkCalculate_VeryLong(b *testing.B) {
	calc := setupCalculator(b)
	input := strings.Repeat("101010", 100) // 600 bits
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_ExtremelyLong benchmarks an extremely long input string.
// This might stress-test the string iteration and FSM transitions.
func BenchmarkCalculate_ExtremelyLong(b *testing.B) {
	calc := setupCalculator(b)
	input := strings.Repeat("110101001", 1000) // 9000 bits
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_EmptyString benchmarks the empty string case.
func BenchmarkCalculate_EmptyString(b *testing.B) {
	calc := setupCalculator(b)
	input := ""
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r, _ := calc.Calculate(input)
		result = r
	}
}

// BenchmarkCalculate_InvalidInput benchmarks an input with an invalid character.
func BenchmarkCalculate_InvalidInput(b *testing.B) {
	calc := setupCalculator(b)
	input := "10101X10101" // Invalid character 'X'
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// For invalid input, an error is expected, but the performance cost of validation
		// is what we are interested in.
		_, _ = calc.Calculate(input) 
	}
}