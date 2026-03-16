//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 상태 다이어그램의 계약 해시를 계산한다
package contract

import (
	"sort"
	"strings"

	"github.com/geul-org/fullend/internal/statemachine"
)

// HashStateDiagram computes a contract hash for a state machine.
// Based on: sorted states + sorted transitions (from:event:to).
func HashStateDiagram(sd *statemachine.StateDiagram) string {
	var parts []string

	states := make([]string, len(sd.States))
	copy(states, sd.States)
	sort.Strings(states)
	parts = append(parts, strings.Join(states, ","))

	var transitions []string
	for _, t := range sd.Transitions {
		transitions = append(transitions, t.From+":"+t.Event+":"+t.To)
	}
	sort.Strings(transitions)
	parts = append(parts, strings.Join(transitions, ","))

	return Hash7(strings.Join(parts, "|"))
}
