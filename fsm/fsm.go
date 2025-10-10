package fsm

import "fmt"

// Automaton decouples consumers from the concrete implementation details
// of the Run method, allowing different FSM types to be plugged in.
type Automaton interface {
	Run(input string) (finalState string, err error)
	IsAccepting(state string) bool
	ValidateInput(input string) bool
}

// FiniteAutomaton (FA) structure
// Represents the 5-tuple: (Q, Σ, q0, F, δ)
type FiniteAutomaton struct {
	States          map[string]bool              // Q: Set of states (S0, S1, S2)
	Alphabet        map[string]bool              // Σ: Input alphabet ('0', '1')
	InitialState    string                       // q0: Initial state (S0)
	AcceptingStates map[string]bool              // F: Set of accepting states (S0, S1, S2 for our use case)
	Transitions     map[string]map[string]string // δ: Transition function: map[CurrentState]map[InputSymbol]NextState
}

// -----------------------------------------------------------------------------
// Generic FSM API Method: Run
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

func NewFiniteAutomaton(
	states []string,
	alphabet []string,
	initialState string,
	acceptingStates []string,
	transitions map[string]map[string]string,
) (Automaton, error) {

	// Create map representations for O(1) lookups
	stateSet := make(map[string]bool)
	for _, s := range states {
		stateSet[s] = true
	}
	alphaSet := make(map[string]bool)
	for _, a := range alphabet {
		alphaSet[a] = true
	}
	acceptingSet := make(map[string]bool)
	for _, f := range acceptingStates {
		acceptingSet[f] = true
	}

	fa := &FiniteAutomaton{
		States:          stateSet,
		Alphabet:        alphaSet,
		InitialState:    initialState,
		AcceptingStates: acceptingSet,
		Transitions:     transitions,
	}

	// 1. Validate Initial State is a member of Q
	if _, ok := stateSet[initialState]; !ok {
		return nil, fmt.Errorf("FSM Config Error: Initial state '%s' is not defined in the set of States (Q)", initialState)
	}

	// 2. Validate Accepting States are a subset of Q
	for _, as := range acceptingStates {
		if _, ok := stateSet[as]; !ok {
			return nil, fmt.Errorf("FSM Config Error: Accepting state '%s' is not defined in the set of States (Q)", as)
		}
	}

	// 3. Validate Transition Completeness (DFA property)
	// Check that every state on every alphabet symbol has a valid transition defined and leads to a valid state.
	for _, currentState := range states {
		transitionsFromCurrent, ok := transitions[currentState]
		if !ok {
			return nil, fmt.Errorf("FSM Config Error: Missing transition rules for state '%s' (not in δ)", currentState)
		}

		for _, symbol := range alphabet {
			nextState, ok := transitionsFromCurrent[symbol]
			if !ok {
				return nil, fmt.Errorf("FSM Config Error: Missing transition for state '%s' on symbol '%s'", currentState, symbol)
			}
			// Check that the resulting nextState is also a member of Q
			if _, ok := stateSet[nextState]; !ok {
				return nil, fmt.Errorf("FSM Config Error: Transition from '%s' on '%s' leads to undefined state '%s'", currentState, symbol, nextState)
			}
		}
	}

	// If all checks pass, return the valid FA
	return fa, nil
}

func (fa *FiniteAutomaton) IsAccepting(state string) bool {
	if _, exists := fa.AcceptingStates[state]; !exists {
		return false
	}
	return true
}

// Iterate over the input string, which naturally iterates over runes (characters).
func (fa *FiniteAutomaton) ValidateInput(input string) bool {
	for _, char := range input {
		// Convert the rune to the Symbol type (string) for map lookup.
		symbol := string(char)

		// Check for the symbol's existence in the alphabet map (O(1) lookup).
		// If the symbol is not found, the `exists` variable will be false.
		if _, exists := fa.Alphabet[symbol]; !exists {
			// If a single character is not in the alphabet, the input is invalid.
			return false
		}
	}

	// If the loop completes, every symbol in the input is valid.
	return true
}
