package mod3

import (
	"strings"
	"testing"
)

// MockAutomaton is retained here for any future isolated component testing, 
// though the public API tests focus on the concrete implementation.
type MockAutomaton struct {
	MockRun func(input string) (finalState string, err error)
}

// Run implements the Automaton interface by calling the mock function
func (m *MockAutomaton) Run(input string) (finalState string, err error) {
	return m.MockRun(input)
}

func GetTestModThreeConfig() FiniteAutomaton {
	return FiniteAutomaton{
		States:          []string{StateS0, StateS1},
		Alphabet:        []string{Symbol0},
		InitialState:    StateS0,
		// All states are accepting in this design, as the final state IS the remainder.
		AcceptingStates: []string{StateS0},
		
		// Transitions (current state -> input symbol -> next state)
		Transitions: map[string]map[string]string{
			StateS0: {Symbol0: StateS1}, 
			StateS1: {Symbol0: StateS1}, 
		},
	}
}
// -----------------------------------------------------------------------------
// 1. UNIT TEST FOR StateToRemainder
// -----------------------------------------------------------------------------

// TestStateToRemainder is retained as a unit test for the private helper method.
func TestStateToRemainder(t *testing.T) {
	// NOTE: We must instantiate a ModThreeCalculator to access its private methods in tests.
	calc, _ := NewModThreeCalculator(GetModThreeConfig()) 
	concreteCalc := calc.(*ModThreeCalculator) // Safely assume it's the concrete type for testing private helpers

	tests := []struct {
		inputState string
		expected   int
	}{
		{StateS0, 0},
		{StateS1, 1},
		{StateS2, 2},
		{"InvalidState", -1}, // Ensures the default unknown state handler works
	}

	for _, tt := range tests {
		t.Run(tt.inputState, func(t *testing.T) {
			actual := concreteCalc.stateToRemainder(tt.inputState) // Call the private method
			if actual != tt.expected {
				t.Errorf("stateToRemainder(%s): got %d, want %d", tt.inputState, actual, tt.expected)
			}
		})
	}
}


// -----------------------------------------------------------------------------
// 1. PUBLIC API CONTRACT TESTS (Ensuring correct setup and calculations)
// -----------------------------------------------------------------------------

// TestNewModThreeCalculator verifies the factory function for success and error paths.
func TestNewModThreeCalculator(t *testing.T) {
	// Sub-test 1: Successful initialization (Covers success path of the factory)
	t.Run("Success", func(t *testing.T) {
		_, err := NewModThreeCalculator(GetModThreeConfig())
		if err != nil {
			t.Fatalf("NewModThreeCalculator failed to initialize FSM engine: %v", err)
		}
	})

	// Sub-test 2: Initialization failure (Covers the hardcoded error return line)
	t.Run("FSM_Init_Failure", func(t *testing.T) {
		// Create an invalid configuration that NewFiniteAutomaton will reject.
		cfg := GetModThreeConfig()
		
		// Invalidate the config by making the transitions empty, which violates DFA rules.
		cfg.Transitions = make(map[string]map[string]string) 
		
		_, err := NewModThreeCalculator(cfg)
		
		// Assert that an error was returned.
		if err == nil {
			t.Fatal("Expected initialization error due to invalid FSM config, but got nil")
		}
		
		// Assert the error message contains the expected factory wrapper text.
		expectedErrSubstring := "failed to initialize FSM engine"
		if !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected factory error containing %q, but got %v", expectedErrSubstring, err)
		}
	})

	// Sub-test 3: Explicit configuration validation failure
	t.Run("Config_Validation_Failure", func(t *testing.T) {
		// Set up an obviously invalid configuration: Initial state is not in the States list.
		cfg := GetModThreeConfig()
		cfg.States = []string{StateS1, StateS2}
		cfg.InitialState = StateS0 // S0 is defined, but not included in States Q
		
		_, err := NewModThreeCalculator(cfg)
		
		if err == nil {
			t.Fatal("Expected initialization error for invalid config, but got nil")
		}
		
		// Assert the error message contains a substring related to the FSM validation.
		expectedErrSubstring := "Initial state"
		if !strings.Contains(err.Error(), expectedErrSubstring) {
			t.Errorf("Expected validation error containing %q (about the initial state), but got %v", expectedErrSubstring, err)
		}
	})
}

// TestCalculator_ErrorPaths verifies that invalid input results in the specified error contract:
// a remainder of -1 and a non-nil error.
func TestCalculator_ErrorPaths(t *testing.T) {
	calc, err := NewModThreeCalculator(GetModThreeConfig())
	if err != nil {
		t.Fatalf("Failed to initialize ModuloCalculator: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected int // Should always be -1 on error
	}{
		{"InvalidSymbol_A", "1A0", -1},
		{"InvalidSymbol_2", "1121", -1},
		{"InvalidSymbol_Space", "1 0", -1},
		{"InvalidSymbol_dot", "1.0", -1},
		{"InvalidSymbol_dash", "1-0", -1},
		{"InvalidSymbol_newline", "1\n0", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, execErr := calc.Calculate(tt.input)

			// 1. Assert that an error occurred
			if execErr == nil {
				t.Fatalf("Calculate(%s) expected an error but got nil (Result: %d)", tt.input, actual)
			}

			// 2. Assert the required error remainder code is returned
			if actual != tt.expected {
				t.Errorf("Calculate(%s) error path failed. Got remainder %d, want %d", tt.input, actual, tt.expected)
			}
		})
	}
}

// TestCalculate_UnknownState explicitly verifies the error path for an FSM result 
// that cannot be mapped to a remainder (0, 1, or 2).
func TestCalculate_UnknownState(t *testing.T) {
    unknownState := "S99"
    
    // Setup a mock Automaton that returns a garbage state (S99) and no FSM error.
    mockFA := &MockAutomaton{
        MockRun: func(input string) (string, error) {
            return unknownState, nil 
        },
    }

	// Manually create the calculator with the mock FA, allowing us to test its internal logic.
	calc := &ModThreeCalculator{fa: mockFA}
	
	actualRemainder, actualErr := calc.Calculate("dummy")

	expectedRemainder := -1
	if actualRemainder != expectedRemainder {
		t.Errorf("Calculate() unknown state failed. Got remainder %d, want %d", actualRemainder, expectedRemainder)
	}
	// Assert the specific "unknown state" error message provided by the user.
	if actualErr == nil || !strings.Contains(actualErr.Error(), "unknown state") {
		t.Errorf("Calculate() error mismatch. Expected 'unknown state' error, got: %v", actualErr)
	}
}

// TestCalculate_AcceptanceCheck verifies the logic inside Calculate for non-accepting states.
func TestCalculate_AcceptanceCheck(t *testing.T) {
    // 1. Configure a concrete FSM that has a valid transition, but terminates in a non-accepting state.
    // States: {S0, S1}. Initial: S0. Accepting: {S0}.
    // Transition: S0 on '0' -> S1. (Input "0" results in non-accepting state S1)
    
    // We must use a concrete FiniteAutomaton to trigger the downcasting check in Calculate.
    fa, faErr := NewFiniteAutomaton(
        []string{StateS0, StateS1}, 
        []string{Symbol0}, 
        StateS0, 
        []string{StateS0}, // Accepting only S0
        map[string]map[string]string{
            StateS0: {Symbol0: StateS1}, // Transition to the non-accepting state
            StateS1: {Symbol0: StateS1}, 
        },
    )

    if faErr != nil {
        t.Fatalf("Failed to create restrictive FiniteAutomaton for test: %v", faErr)
    }
    
    // 2. Initialize the calculator with the restrictive FA.
    calc := &ModThreeCalculator{fa: fa}
	
    // 3. Call Calculate with an input ("0") that forces the FSM to transition to the non-accepting state (S1).
	actualRemainder, actualErr := calc.Calculate("0") 

	expectedRemainder := -1
	if actualRemainder != expectedRemainder {
		t.Errorf("Calculate() acceptance check failed. Got remainder %d, want %d", actualRemainder, expectedRemainder)
	}
	if actualErr == nil || !strings.Contains(actualErr.Error(), "non-accepting state") {
		t.Errorf("Calculate() error mismatch. Expected 'non-accepting state' error, got: %v", actualErr)
	}
}


// TestCalculator_Correctness tests the public contract (Calculate) with valid binary inputs,
// including edge cases and large numbers.
func TestCalculator_Correctness(t *testing.T) {
	// Initialize the calculator service once for all subsequent tests
	calc, err := NewModThreeCalculator(GetModThreeConfig())
	if err != nil {
		t.Fatalf("Failed to initialize ModuloCalculator: %v", err)
	}

	thirtyTwoZeros := strings.Repeat("0", 32)
	thirtyTwoOnes := strings.Repeat("1", 32)
	sixtyFourZeros := strings.Repeat("0", 64)
	sixtyFourOnes := strings.Repeat("1", 64)

	tests := []struct {
		name     string
		input    string // Binary input
		expected int    // Expected remainder
	}{
		// Edge Cases
		{"EmptyString", "", 0},
		{"LeadingZeros", "000101", 2}, // 5 mod 3 = 2
		{"SingleZero", "0", 0},
		{"SingleOne", "1", 1},

		// Small Value Tests
		{"Zero", "0", 0},
		{"One", "1", 1},
		{"Two", "10", 2},      // 2 mod 3 = 2
		{"Three", "11", 0},     // 3 mod 3 = 0
		{"Four", "100", 1},    // 4 mod 3 = 1
		{"Five", "101", 2},    // 5 mod 3 = 2
		{"Six", "110", 0},     // 6 mod 3 = 0
		{"Seven", "111", 1},    // 7 mod 3 = 1
		{"Eight", "1000", 2},   // 8 mod 3 = 2
		{"Nine", "1001", 0},    // 9 mod 3 = 0
		{"Ten", "1010", 1},     // 10 mod 3 = 1
		{"Binary341", "101010101", 2}, // 341 mod 3 = 2
		{"Binary42", "101010", 0},     // 42 mod 3 = 0
		{"2Power8", "1" + strings.Repeat("0", 8), 1}, // 256 mod 3 = 1

		// Large Value Tests
		{"2Power32", "1" + thirtyTwoZeros, 1},       // 2^32 mod 3 = 1
		{"32Ones", thirtyTwoOnes, 0},               // (2^32 - 1) mod 3 = 0 (since 2^32 mod 3 = 1, then 1-1 = 0)
		{"2Power64", "1" + sixtyFourZeros, 1},       // 2^64 mod 3 = 1
		{"64Ones", sixtyFourOnes, 0},               // (2^64 - 1) mod 3 = 0

		// Alternating patterns
		{"Alternating10_Short", "1010", 1},         // 10 (binary) = 2, 1010 (binary) = 10, 10 mod 3 = 1
		{"Alternating01_Short", "0101", 2},         // 0101 (binary) = 5, 5 mod 3 = 2 (initial zeros are effectively ignored)
		{"Alternating10_Long", strings.Repeat("10", 30), 0}, // 60 bits, should result in 0

		// Powers of 2 close to 3
		{"2Power1", "10", 2},
		{"2Power2", "100", 1},
		{"2Power3", "1000", 2},
		{"2Power4", "10000", 1},
		{"2Power5", "100000", 2},
		{"2Power6", "1000000", 1},

		// Large Value Tests
		{"2Power32", "1" + thirtyTwoZeros, 1},
		{"32Ones", thirtyTwoOnes, 0},

		// Stress Test: 100 repetitions of "10" (200 bits total)
		{"Huge200Bits", strings.Repeat("10", 100), 2},
		{"LongStringOfOnes", strings.Repeat("1", 20), 0}, // (2^20 - 1) mod 3 = 0
		{"LongStringOfZerosWithTrailingOne", strings.Repeat("0", 15) + "1", 1}, // 1 mod 3 = 1
		{"OneFollowedByManyZeros", "1" + strings.Repeat("0", 19), 2}, // 2^19 mod 3 = 2 

		// More mixed patterns
		{"MixedPattern1", "110101101", 0}, // 429 mod 3 = 0
		{"MixedPattern2", "100110111", 2}, // 311 mod 3 = 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the public interface method
			actual, execErr := calc.Calculate(tt.input)

			if execErr != nil {
				t.Errorf("Calculate(%s) failed unexpectedly with error: %v", tt.input, execErr)
				return
			}
			if actual != tt.expected {
				t.Errorf("Calculate(%s): got remainder %d, want %d", tt.input, actual, tt.expected)
			}
		})
	}
}