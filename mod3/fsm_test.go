package mod3

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
	states := []string{StateS0, StateS1, StateS2}
	alphabet := []string{Symbol0, Symbol1}
	initialState := StateS0
	acceptingStates := []string{StateS0, StateS1, StateS2}
	transitions := map[string]map[string]string{
		StateS0: {Symbol0: StateS0, Symbol1: StateS1},
		StateS1: {Symbol0: StateS2, Symbol1: StateS0},
		StateS2: {Symbol0: StateS1, Symbol1: StateS2},
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
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: []string{StateS0, "S99"}, transitions: transitions,
			expectError: true, errorContains: "Accepting state 'S99' is not defined",
		},
		{
			name: "Error: Missing Transition Rules (State S0 missing)",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{StateS1: transitions[StateS1], StateS2: transitions[StateS2]}, // S0 removed
			expectError: true, errorContains: "Missing transition rules for state 'S0'",
		},
		{
			name: "Error: Missing Transition for Symbol (S0 missing '1')",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{
				StateS0: {Symbol0: StateS0}, // Symbol1 is missing
				StateS1: transitions[StateS1], 
				StateS2: transitions[StateS2],
			},
			expectError: true, errorContains: "Missing transition for state 'S0' on symbol '1'",
		},
		{
			name: "Error: Transition Leads to Undefined State",
			states: states, alphabet: alphabet, initialState: initialState, acceptingStates: acceptingStates, 
			transitions: map[string]map[string]string{
				StateS0: {Symbol0: "S99", Symbol1: StateS1}, // S0 on 0 -> S99 (bad state)
				StateS1: transitions[StateS1], 
				StateS2: transitions[StateS2],
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