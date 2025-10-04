
package mod3

import (
	"strings"
	"fmt"
)

const (
	// Define the States
    StateS0 = "S0"
    StateS1 = "S1"
    StateS2 = "S2"
    
    // Define input symbols too
    Symbol0 = "0"
    Symbol1 = "1"
)

type ModuloCalculator interface {
    Calculate(input string) (remainder int, err error)
}

type ModThreeCalculator struct {
    fa Automaton // The underlying generic FSM engine.
}

func GetModThreeConfig() FiniteAutomaton {
	return FiniteAutomaton{
		States:          []string{StateS0, StateS1, StateS2},
		Alphabet:        []string{Symbol0, Symbol1},
		InitialState:    StateS0,
		// All states are accepting in this design, as the final state IS the remainder.
		AcceptingStates: []string{StateS0, StateS1, StateS2},
		
		// Transitions (current state -> input symbol -> next state)
		Transitions: map[string]map[string]string{
			StateS0: {Symbol0: StateS0, Symbol1: StateS1}, 
			StateS1: {Symbol0: StateS2, Symbol1: StateS0}, 
			StateS2: {Symbol0: StateS1, Symbol1: StateS2}, 
		},
	}
}

// NewModThreeCalculator initializes the calculator using the separated configuration.
func NewModThreeCalculator(cfg FiniteAutomaton) (ModuloCalculator, error) {
	// Pass the structured configuration data to the FSM constructor
	// Here is better to passing FiniteAutomaton for initialization to make it more loosely coupled
	fa, err := NewFiniteAutomaton(cfg.States, cfg.Alphabet, cfg.InitialState, cfg.AcceptingStates, cfg.Transitions)
	
	// This is the error path you wanted to ensure is covered.
	if err != nil {
		// This line will now only execute if GetModThreeConfig() contains an invalid definition.
		return nil, fmt.Errorf("failed to initialize FSM engine: %w", err)
	}
	
	return &ModThreeCalculator{fa: fa}, nil
}

// --- PRIVATE HELPER METHODS ---

// stateToRemainder maps the final state to the required remainder (0, 1, or 2).
func (c *ModThreeCalculator) stateToRemainder(state string) int {
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

// isStateAccepting checks if the final state is one of the designated accepting states.
func (c *ModThreeCalculator) isStateAccepting(finalState string) bool {
    return c.fa.IsAccepting(finalState) // No type assertion needed!
}

// --- PUBLIC INTERFACE METHOD IMPLEMENTATION ---

// Calculate runs the binary input through the configured FSM and returns the final remainder.
// This implements the ModuloCalculator interface.
func (c *ModThreeCalculator) Calculate(input string) (int, error) {
	// Handle empty string case (value 0, remainder 0)
	if strings.TrimSpace(input) == "" {
		return 0, nil
	}

	// 1. Run the input against the generic FA engine
	finalState, err := c.fa.Run(input)
	if err != nil {
		return -1, err
	}

	// 2. Acceptance Check
	if !c.isStateAccepting(finalState) {
		return -1, fmt.Errorf("FSM execution ended in non-accepting state: %s", finalState)
	}

	// 3. Map the resulting state to the remainder output
	remainder := c.stateToRemainder(finalState)
	if remainder == -1 {
		// Should only happen if finalState is totally unexpected (e.g. "S99")
		return -1, fmt.Errorf("FSM execution resulted in unknown state: %s", finalState)
	}
	
	return remainder, nil
}
