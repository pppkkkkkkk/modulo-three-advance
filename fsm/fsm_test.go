package fsm

import (
	"testing"
	"errors"
	"strings"
)

// setupSimpleFA creates a minimal FSM configuration for testing the Run method's core logic.
// This FA moves from Start -> Middle on 'a', and Middle -> End on 'b'.
func setupSimpleFA() *FiniteAutomaton {
	return &FiniteAutomaton{
		InitialState: "Start",
		Transitions: map[string]map[string]string{
			"Start": {
				"a": "Middle", // Transition 1: Start -> Middle on 'a'
				"x": "Fail",   // Allow 'x' but lead to a non-final/trap state for testing
			},
			"Middle": {
				"b": "End",    // Transition 2: Middle -> End on 'b'
			},
			"End": {
				"c": "End",    // Transition 3: Loop at End on 'c'
			},
		},
	}
}

// -----------------------------------------------------------------------------
// 1. UNIT TEST FOR FiniteAutomaton Run
// -----------------------------------------------------------------------------

func TestFiniteAutomaton_Run(t *testing.T) {
	fa := setupSimpleFA()

	tests := []struct {
		name         string
		input        string
		expectedState string
		expectedError error
	}{
		// 1. Successful Paths
		{"EmptyInput", "", "Start", nil},                             // Should return InitialState
		{"SingleStep", "a", "Middle", nil},                           // Start -> Middle
		{"TwoSteps", "ab", "End", nil},                               // Start -> Middle -> End
		{"LoopAtEnd", "abc", "End", nil},                             // Start -> Middle -> End -> End

		// 2. Invalid Input Symbol Handling (Errors)
		{"InvalidSymbol_Start", "b", "", errors.New("FSM Error: Invalid input symbol 'b' for state Start")}, // 'b' not valid in 'Start'
		{"InvalidSymbol_Middle", "aab", "", errors.New("FSM Error: Invalid input symbol 'a' for state Middle")}, // 'a' not valid in 'Middle' after 'a'

		// 3. Missing State/Transition Handling (Errors)
		// and then trigger the "Transition rule missing" error for the subsequent symbol 'b'.
		{"MissingTransition", "xb", "", errors.New("FSM Error: Transition rule missing for state Fail")}, 
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualState, actualErr := fa.Run(tt.input)

			// Check for state match
			if actualState != tt.expectedState {
				t.Errorf("Run(%q) state mismatch. Got %q, want %q", tt.input, actualState, tt.expectedState)
			}

			// Check for error match (allowing nil vs nil, or specific error text)
			if tt.expectedError == nil {
				if actualErr != nil {
					t.Errorf("Run(%q) expected no error, but got: %v", tt.input, actualErr)
				}
			} else {
				// Only compare error string if the error exists
				if actualErr == nil || actualErr.Error() != tt.expectedError.Error() {
					t.Errorf("Run(%q) error mismatch. Got %v, want %v", tt.input, actualErr, tt.expectedError)
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 2. UNIT TEST FOR NewFiniteAutomaton
// -----------------------------------------------------------------------------

func TestNewFiniteAutomaton_Validation(t *testing.T) {
	// Base valid configuration (Mod 3 FSM)
	states := []string{"S0", "S1", "S2"}
	alphabet := []string{"0", "1"}
	initialState := "S0"
	acceptingStates := []string{"S0", "S1", "S2"}
	transitions := map[string]map[string]string{
		"S0": {"0": "S0", "1": "S1"},
		"S1": {"0": "S2", "1": "S0"},
		"S2": {"0": "S1", "1": "S2"},
	}
	
	tests := []struct {
		name string
		states []string
		alphabet []string
		initialState string
		acceptingStates []string
		transitions map[string]map[string]string
		expectError bool
		errorContains string
	}{
		{
			name: "Valid Modulo 3 Config (Success)",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, transitions: transitions,
			expectError: false,
		},
		{
			name: "Error: Initial State Not in Q",
			states: states, alphabet: alphabet, initialState: "S99", acceptingStates: acceptingStates, transitions: transitions,
			expectError: true, errorContains: "Initial state 'S99' is not defined",
		},
		{
			name: "Error: Accepting State Not in Q",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: []string{"S0", "S99"}, transitions: transitions,
			expectError: true, errorContains: "Accepting state 'S99' is not defined",
		},
		{
			name: "Error: Missing Transition Rules (State S0 missing)",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{"S1": transitions["S1"], "S2": transitions["S2"]}, // S0 removed
			expectError: true, errorContains: "Missing transition rules for state 'S0'",
		},
		{
			name: "Error: Missing Transition for Symbol (S0 missing '1')",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{
				"S0": {"0": "S0"}, // "1" is missing
				"S1": transitions["S1"], 
				"S2": transitions["S2"],
			},
			expectError: true, errorContains: "Missing transition for state 'S0' on symbol '1'",
		},
		{
			name: "Error: Transition Leads to Undefined State",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{
				"S0": {"0": "S99", "1": "S1"}, // S0 on 0 -> S99 (bad state)
				"S1": transitions["S1"], 
				"S2": transitions["S2"],
			},
			expectError: true, errorContains: "leads to undefined state 'S99'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFiniteAutomaton(tt.states, tt.alphabet, tt.initialState, tt.acceptingStates, tt.transitions)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing %q, but got nil", tt.errorContains)
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error content mismatch. Got error: %q, but expected it to contain: %q", err.Error(), tt.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 3. UNIT TEST FOR IsAccepting
// -----------------------------------------------------------------------------

// TestFiniteAutomaton_IsAccepting tests the IsAccepting method of the FiniteAutomaton
func TestFiniteAutomaton_IsAccepting(t *testing.T) {
	tests := []struct {
		name            string
		acceptingStates map[string]bool // The F: Set of accepting states for the FA
		stateToCheck    string          // The state whose acceptance status is being checked
		expected        bool            // Expected result of IsAccepting
	}{
		// --- Test cases where the state *is* accepting ---
		{"IsAccepting_ExistingState_IsAccepting",
			map[string]bool{"S0": true, "S1": true, "S2": true}, // All states accepting
			"S1",
			true},
		{"IsAccepting_OnlyOneAccepting_IsAccepting",
			map[string]bool{"S_final": true}, // Only S_final is accepting
			"S_final",
			true},
		{"IsAccepting_InitialStateIsAccepting",
			map[string]bool{"Q_start": true, "Q_mid": true}, // Q_start is accepting
			"Q_start",
			true},

		// --- Test cases where the state *is not* accepting ---
		{"IsAccepting_ExistingState_NotAccepting",
			map[string]bool{"S0": true, "S2": true}, // S1 is not listed
			"S1",
			false},
		{"IsAccepting_NonExistentState",
			map[string]bool{"S0": true, "S1": true}, // "S_unknown" is not in any FSM config
			"S_unknown",
			false},
		{"IsAccepting_EmptyAcceptingSet",
			map[string]bool{}, // No accepting states
			"S0",
			false},
		{"IsAccepting_EmptyAcceptingSet_EmptyState",
			map[string]bool{}, // No accepting states
			"", // An empty string state (unlikely but possible)
			false},
		{"IsAccepting_StateNotInSetButEmptySet",
			map[string]bool{"A": true},
			"B",
			false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy FiniteAutomaton just to hold the AcceptingStates for this test.
			// The other fields are irrelevant for IsAccepting.
			fa := &FiniteAutomaton{
				AcceptingStates: tt.acceptingStates, // Use the test-specific accepting states map
			}

			actual := fa.IsAccepting(tt.stateToCheck)

			if actual != tt.expected {
				t.Errorf("IsAccepting(%q) with accepting states %v: got %t, want %t",
					tt.stateToCheck, tt.acceptingStates, actual, tt.expected)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 4. UNIT TEST FOR ValidateInput
// -----------------------------------------------------------------------------

func TestFiniteAutomaton_ValidateInput(t *testing.T) {
	// Define a common alphabet for testing
	commonAlphabet := map[string]bool{
		"a": true,
		"b": true,
		"c": true,
		"0": true,
		"1": true,
	}

	tests := []struct {
		name     string
		alphabet map[string]bool
		input    string
		expected bool
	}{
		{
			name:     "Valid input with all symbols in alphabet",
			alphabet: commonAlphabet,
			input:    "abc10",
			expected: true,
		},
		{
			name:     "Valid input with a single symbol",
			alphabet: commonAlphabet,
			input:    "a",
			expected: true,
		},
		{
			name:     "Empty input",
			alphabet: commonAlphabet,
			input:    "",
			expected: true, // An empty string contains no invalid characters
		},
		{
			name:     "Input with one invalid symbol",
			alphabet: commonAlphabet,
			input:    "abcz",
			expected: false,
		},
		{
			name:     "Input with multiple invalid symbols",
			alphabet: commonAlphabet,
			input:    "xyz",
			expected: false,
		},
		{
			name:     "Input with a mix of valid and invalid symbols (invalid first)",
			alphabet: commonAlphabet,
			input:    "zabc",
			expected: false,
		},
		{
			name:     "Input with a mix of valid and invalid symbols (invalid in middle)",
			alphabet: commonAlphabet,
			input:    "abz01",
			expected: false,
		},
		{
			name:     "Alphabet with only one symbol, valid input",
			alphabet: map[string]bool{"x": true},
			input:    "xxx",
			expected: true,
		},
		{
			name:     "Alphabet with only one symbol, invalid input",
			alphabet: map[string]bool{"x": true},
			input:    "xxy",
			expected: false,
		},
		{
			name:     "Empty alphabet, valid input (empty string)",
			alphabet: map[string]bool{},
			input:    "",
			expected: true,
		},
		{
			name:     "Empty alphabet, invalid input (non-empty string)",
			alphabet: map[string]bool{},
			input:    "a",
			expected: false,
		},
		{
			name:     "Input with special characters (valid)",
			alphabet: map[string]bool{"@": true, "#": true, "$": true},
			input:    "@#$",
			expected: true,
		},
		{
			name:     "Input with special characters (invalid)",
			alphabet: map[string]bool{"@": true},
			input:    "@#",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fa := &FiniteAutomaton{
				Alphabet: tt.alphabet,
			}
			actual := fa.ValidateInput(tt.input)
			if actual != tt.expected {
				t.Errorf("ValidateInput() for input '%s' with alphabet %v got %v, want %v", tt.input, tt.alphabet, actual, tt.expected)
			}
		})
	}
}