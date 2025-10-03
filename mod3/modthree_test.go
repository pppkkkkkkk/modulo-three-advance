package mod3

import (
	"testing"
	"strings"
	"errors"
)

// MockAutomaton is a mock implementation of the Automaton interface
type MockAutomaton struct {
	MockRun func(input string) (finalState string, err error)
}

// Run implements the Automaton interface by calling the mock function
func (m *MockAutomaton) Run(input string) (finalState string, err error) {
	return m.MockRun(input)
}

// -----------------------------------------------------------------------------
// 1. UNIT TEST FOR ModThreeFA (Constructor) & StateToRemainder (Mapper)
// -----------------------------------------------------------------------------

func TestModThreeFA_Initialization(t *testing.T) {
	// Use type assertion to access the concrete FiniteAutomaton fields
	fa := ModThreeFA().(*FiniteAutomaton) 

	// 1. Check Initial State (q0)
	if fa.InitialState != StateS0 { 
		t.Errorf("InitialState mismatch. Got %s, want %s", fa.InitialState, StateS0)
	}

	// 2. Check States (Q)
	if len(fa.States) != 3 {
		t.Errorf("States count mismatch. Got %d, want 3 (S0, S1, S2)", len(fa.States))
	}

	// 3. Check Alphabet (Σ)
	if len(fa.Alphabet) != 2 {
		t.Errorf("Alphabet count mismatch. Got %d, want 2 (0, 1)", len(fa.Alphabet))
	}

	// 4. Check Transitions (δ) Structure size
	if len(fa.Transitions) != 3 {
		t.Fatalf("Transitions map size mismatch. Got %d entries, want 3 (for S0, S1, S2)", len(fa.Transitions))
	}
	
	// 5. Check a Representative Transition (S1 on '1' -> S0)
	// This ensures the core logic mapping is present and correct: (1 * 2 + 1) mod 3 = 0.
	expectedNextState := StateS0
	actualNextState, ok := fa.Transitions[StateS1][Symbol1]
	
	if !ok {
		t.Errorf("Transition rule missing for state %s on symbol %s", StateS1, Symbol1)
	} else if actualNextState != expectedNextState {
		t.Errorf("Transition S1 on '1' failed. Got %s, want %s", actualNextState, expectedNextState)
	}
}

func TestStateToRemainder(t *testing.T) {
	tests := []struct {
		inputState string
		expected   int
	}{
		{StateS0, 0}, 
		{StateS1, 1}, 
		{StateS2, 2},
		{"InvalidState", -1}, 
	}

	for _, tt := range tests {
		t.Run(tt.inputState, func(t *testing.T) {
			actual := StateToRemainder(tt.inputState) 
			if actual != tt.expected {
				t.Errorf("StateToRemainder(%s): got %d, want %d", tt.inputState, actual, tt.expected)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 2. ISOLATED UNIT TEST FOR modThree CORE LOGIC (Using Mocking)
// -----------------------------------------------------------------------------

func TestModThree_Isolated_ErrorPath(t *testing.T) {
	expectedError := errors.New("Simulated FSM Failure")

	// 1. Setup Mock FA configured to fail on Run()
	mockFA := &MockAutomaton{
		MockRun: func(input string) (string, error) {
			return "", expectedError // Force Run to return an error
		},
	}

	// 2. Call the core logic function directly with the mock
	actualRemainder := modThree(mockFA, "101") 

	// 3. Assert that the error path returned the expected remainder of -1
	expectedRemainder := -1
	if actualRemainder != expectedRemainder {
		t.Errorf("modThree() error path failed. Got %d, want %d", actualRemainder, expectedRemainder)
	}
}

func TestModThree_Isolated_SuccessPath(t *testing.T) {
	// Test the core success path: finalState -> StateToRemainder(finalState)
	// We want to simulate the FA ending in StateS2 (remainder 2)
	expectedState := StateS2
	expectedRemainder := 2 // StateS2 maps to 2

	// 1. Setup Mock FA configured to successfully return StateS2
	mockFA := &MockAutomaton{
		MockRun: func(input string) (string, error) {
			return expectedState, nil // Force Run to return StateS2 and no error
		},
	}

	// 2. Call the core logic function directly with the mock
	actualRemainder := modThree(mockFA, "101") 

	// 3. Assert that the result maps correctly
	if actualRemainder != expectedRemainder {
		t.Errorf("modThree() success path failed. Got remainder %d, want %d (mapped from state %s)", actualRemainder, expectedRemainder, expectedState)
	}
}

func TestModThree_Isolated_WhitespaceInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty String", "", 0},
		{"Single Space", " ", 0},
		{"Tabs and Newlines", "\t\n  \r", 0},
	}

	// NOTE: We don't even need a fully configured mock here, as the function should return 0 before calling fa.Run.
	// We use the simplest possible mock just to satisfy the function signature.
	mockFA := &MockAutomaton{} 

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := modThree(mockFA, tt.input)
			if actual != tt.expected {
				t.Errorf("modThree(%q) whitespace check failed. Got %d, want %d", tt.input, actual, tt.expected)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 3. INTEGRATION TEST FOR ModThree (Public API)
// -----------------------------------------------------------------------------

func TestModThree_Correctness(t *testing.T) {
	thirtyTwoZeros := strings.Repeat("0", 32)
	thirtyTwoOnes := strings.Repeat("1", 32)
	
	tests := []struct {
		input    string // Binary input
		expected int    // Expected remainder
		comment  string // Required third field
	}{
		// Edge Cases
		{"", 0, "Empty string"},
		{"000101", 2, "Leading zeros (5 mod 3 = 2)"},
		
		// Small Value Tests
		{"1", 1, "1 mod 3 = 1"}, 
		{"110", 0, "6 mod 3 = 0"},
		{"101010101", 2, "Binary 341 mod 3 = 2"},

		// --- New Medium/Small Tests ---
        {"101010", 0, "Binary 42 mod 3 = 0"}, 
        {"1010011", 2, "Binary 83 mod 3 = 2"}, // Test case resulting in remainder 2
        {"1" + strings.Repeat("0", 8), 1, "2^8 (256) mod 3 = 1"}, // Test power of 2 (even exponent)

		// Large Value Tests
		{"1" + thirtyTwoZeros, 1, "2^32 mod 3 = 1"}, 
		{thirtyTwoOnes, 0, "2^32 - 1 mod 3 = 0"}, 

		// Huge String Test: 100 repetitions of "10" (200 bits total)
		{strings.Repeat("10", 100), 2, "Huge 200-bit string mod 3 = 2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// This calls the public ModThree(input string) wrapper, testing the entire pipeline.
			actual := ModThree(tt.input)
			if actual != tt.expected {
				t.Errorf("ModThree(%s) [%s]: got remainder %d, want %d", tt.input, tt.comment, actual, tt.expected)
			}
		})
	}
}
