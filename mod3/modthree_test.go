package mod3

import (
	"strings"
	"testing"
	"errors"
	"modulo_three_advanced/fsm" 
)

// MockAutomaton is retained here for any future isolated component testing, 
// though the public API tests focus on the concrete implementation.
type MockAutomaton struct {
	MockRun func(input string) (finalState string, err error)
	MockIsAccepting func(state string) bool 
	MockValidateInput	func(input string) bool
}

// Use the below mock functions for Automaton testing 
func (m *MockAutomaton) Run(input string) (finalState string, err error) {
	return m.MockRun(input)
}
func (m *MockAutomaton) IsAccepting(input string) bool {
	return m.MockIsAccepting(input)
}
func (m *MockAutomaton) ValidateInput(input string) bool {
	return m.MockValidateInput(input)
}

// -----------------------------------------------------------------------------
// 1. UNIT TEST FOR StateToRemainder and IsStateAccepting
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

func TestIsStateAccepting(t *testing.T) {
    mockFA := &MockAutomaton{
        MockIsAccepting: func(state string) bool {
            return state == StateS0 // Only S0 is accepting for this mock
        },
    }
    calc := &ModThreeCalculator{fa: mockFA}

    tests := []struct {
        state    string
        expected bool
    }{
        {StateS0, true},
        {StateS1, false},
        {StateS2, false},
        {"Unknown", false},
    }

    for _, tt := range tests {
        t.Run(tt.state, func(t *testing.T) {
            actual := calc.isStateAccepting(tt.state)
            if actual != tt.expected {
                t.Errorf("isStateAccepting(%s): got %t, want %t", tt.state, actual, tt.expected)
            }
        })
    }
}

// -----------------------------------------------------------------------------
// 2. PUBLIC API CONTRACT TESTS (Ensuring correct setup and 100% unit test coverage)
// -----------------------------------------------------------------------------

// TestNewModThreeCalculator verifies the factory function for success and error paths.
func TestNewModThreeCalculator(t *testing.T) {
	// Sub-test 1: Successful initialization (Covers success path of the factory)
	t.Run("Success", func(t *testing.T) {
		calc, err := NewModThreeCalculator(GetModThreeConfig())
		if err != nil {
			t.Fatalf("NewModThreeCalculator failed to initialize FSM engine: %v", err)
		}
		if calc == nil {
			t.Fatal("NewModThreeCalculator returned a nil calculator, but no error")
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

func TestCalculate_ErrorHandling(t *testing.T) {
    // Test case 1: Input validation failure
    t.Run("InputValidationFailed", func(t *testing.T) {
        calc, _ := NewModThreeCalculator(GetModThreeConfig()) // Valid FSM init
        _, err := calc.Calculate("1A0")
        if err == nil || !strings.Contains(err.Error(), "validate Input") {
            t.Errorf("Expected 'validate Input' error, got %v", err)
        }
    })

    // Test case 2: FSM.Run returns an error (using a mock)
    t.Run("FSMRunError", func(t *testing.T) {
        mockFA := &MockAutomaton{
            MockRun: func(input string) (string, error) { return "", errors.New("mock FSM run error") },
            MockValidateInput: func(input string) bool { return true },
            MockIsAccepting: func(input string) bool { return true }, // Irrelevant for this path
        }
        calc := &ModThreeCalculator{fa: mockFA}
        _, err := calc.Calculate("101")
        if err == nil || !strings.Contains(err.Error(), "mock FSM run error") {
            t.Errorf("Expected 'mock FSM run error', got %v", err)
        }
    })

    // Test case 3: Non-accepting state (using a mock or a specially configured real FSM)
    t.Run("NonAcceptingFinalState", func(t *testing.T) {
        fa, faErr := fsm.NewFiniteAutomaton(
            []string{StateS0, StateS1}, 
			[]string{Symbol0}, 
			StateS0, 
			[]string{StateS0},
            map[string]map[string]string{StateS0: {Symbol0: StateS1}, StateS1: {Symbol0: StateS0}},
        )
		if faErr != nil { 
            t.Fatalf("Failed to create restrictive FiniteAutomaton for test: %v", faErr)
        }
        calc := &ModThreeCalculator{fa: fa}
        _, err := calc.Calculate("0") // Goes to S1, which is not accepting
        if err == nil || !strings.Contains(err.Error(), "non-accepting state") {
            t.Errorf("Expected 'non-accepting state' error, got %v", err)
        }
    })

    // Test case 4: Unknown final state (using a mock)
    t.Run("UnknownFinalState", func(t *testing.T) {
        mockFA := &MockAutomaton{
            MockRun: func(input string) (string, error) { return "S99", nil }, // Unknown state
            MockValidateInput: func(input string) bool { return true },
            MockIsAccepting: func(input string) bool { return true },
        }
        calc := &ModThreeCalculator{fa: mockFA}
        _, err := calc.Calculate("101")
        if err == nil || !strings.Contains(err.Error(), "unknown state") {
            t.Errorf("Expected 'unknown state' error, got %v", err)
        }
    })
}

// -----------------------------------------------------------------------------
// INTEGRATION TEST FOR ModThree (Public API)
// -----------------------------------------------------------------------------

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