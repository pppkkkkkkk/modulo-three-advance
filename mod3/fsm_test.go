package mod3

import (
	"testing"
	"errors"
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
