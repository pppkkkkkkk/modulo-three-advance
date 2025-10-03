
package mod3

import (
	"strings"
)

const (
    StateS0 = "S0"
    StateS1 = "S1"
    StateS2 = "S2"
    
    // Define input symbols too
    Symbol0 = "0"
    Symbol1 = "1"
)

// ModThreeFA returns a fully configured FiniteAutomaton for the modulo-three problem.
func ModThreeFA() Automaton {
	fa := &FiniteAutomaton{
		// Q: States
		States: []string{StateS0, StateS1, StateS2},
		// Σ: Alphabet
		Alphabet: []string{Symbol0, Symbol1},
		// q0: Initial State
		InitialState: StateS0,
		// F: Accepting/Final States (All states are final for remainder determination)
		AcceptingStates: []string{StateS0, StateS1, StateS2},
		// δ: Transition function
		Transitions: map[string]map[string]string{
			StateS0: { // Remainder 0
				Symbol0: StateS0, // (0 * 2 + 0) mod 3 = 0
				Symbol1: StateS1, // (0 * 2 + 1) mod 3 = 1
			},
			StateS1: { // Remainder 1
				Symbol0: StateS2, // (1 * 2 + 0) mod 3 = 2
				Symbol1: StateS0, // (1 * 2 + 1) mod 3 = 0
			},
			StateS2: { // Remainder 2
				Symbol0: StateS1, // (2 * 2 + 0) mod 3 = 1
				Symbol1: StateS2, // (2 * 2 + 1) mod 3 = 2
			},
		},
	}
	return fa
}

// StateToRemainder maps the final state to the required remainder (0, 1, or 2).
func StateToRemainder(state string) int {
	switch state {
	case StateS0:
		return 0
	case StateS1:
		return 1
	case StateS2:
		return 2
	default:
		// Should not happen with valid FSM execution
		return -1
	}
}

func ModThree(input string) int {
	// The core logic function modThree (below) is called with the REAL Finite Automaton.
	return modThree(ModThreeFA(), input)
}

// -----------------------------------------------------------------------------
// 4. The Final modThree Function (Using the Generic FA)
// -----------------------------------------------------------------------------

// modThree implements the required procedure using the configured FA.
func modThree(fa Automaton, input string) int {
	// Handle empty string case (value 0, remainder 0)
	if strings.TrimSpace(input) == "" {
		return 0
	}

	// Run the input against the generic FA engine
	finalState, err := fa.Run(input)
	if err != nil {
		// fmt.Printf("Execution Error: %v\n", err)
		return -1
	}

	// Map the resulting state to the remainder output
	return StateToRemainder(finalState)
}
