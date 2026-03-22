//ff:func feature=gen-gogin type=generator control=sequence topic=states
//ff:what string 필드용 CanTransition 함수 코드를 생성한다

package gogin

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// writeStringCanTransition writes the string-based CanTransition function body.
func writeStringCanTransition(buf *strings.Builder, d *statemachine.StateDiagram) {
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
