package mod3

import "fmt"

// Automaton decouples consumers from the concrete implementation details 
// of the Run method, allowing different FSM types to be plugged in.
type Automaton interface {
    Run(input string) (finalState string, err error)
}

// FiniteAutomaton (FA) structure
// Represents the 5-tuple: (Q, Σ, q0, F, δ)
type FiniteAutomaton struct {
	States           []string                   // Q: Set of states (S0, S1, S2)
	Alphabet         []string                   // Σ: Input alphabet ('0', '1')
	InitialState     string                     // q0: Initial state (S0)
	AcceptingStates  []string                   // F: Set of accepting states (S0, S1, S2 for our use case)
	Transitions map[string]map[string]string	// δ: Transition function: map[CurrentState]map[InputSymbol]NextState
}

// -----------------------------------------------------------------------------
// 2. Generic FSM API Method: Run
// -----------------------------------------------------------------------------

// Run processes an input string against the FA configuration and returns the final state.
func (fa *FiniteAutomaton) Run(input string) (finalState string, err error) {
	// Start at the initial state
	currentState := fa.InitialState

	for _, char := range input {
		symbol := string(char)

		// 1. Check if the current state exists in the transition map
		transitionsFromCurrent, ok := fa.Transitions[currentState]
		if !ok {
			return "", fmt.Errorf("FSM Error: Transition rule missing for state %s", currentState)
		}

		// 2. Check if the input symbol is valid for the current state
		nextState, ok := transitionsFromCurrent[symbol]
		if !ok {
			return "", fmt.Errorf("FSM Error: Invalid input symbol '%s' for state %s", symbol, currentState)
		}

		// 3. Move to the next state
		currentState = nextState
	}

	// The state after the entire string is processed is the final state.
	return currentState, nil
}
