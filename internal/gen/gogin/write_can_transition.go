//ff:func feature=gen-gogin type=generator control=selection
//ff:what writes the CanTransition function based on field type

package gogin

import (
	"strings"

	"github.com/geul-org/fullend/internal/statemachine"
)

// writeCanTransition writes the CanTransition function appropriate for the field type.
func writeCanTransition(buf *strings.Builder, d *statemachine.StateDiagram, fieldType string) {
	switch fieldType {
	case "bool":
		buf.WriteString(generateBoolCanTransition(d))
	default:
		buf.WriteString(`// CanTransition checks if the given event is valid from the current state.
// Set DISABLE_STATE_CHECK=1 to bypass state transition checks.
func CanTransition(input Input, event string) error {
	if os.Getenv("DISABLE_STATE_CHECK") == "1" {
		return nil
	}
	status, _ := input.Status.(string)
	_, ok := transitions[transitionKey{from: status, event: event}]
	if !ok {
		return fmt.Errorf("cannot transition from %q via %q", status, event)
	}
	return nil
}
`)
	}
}
