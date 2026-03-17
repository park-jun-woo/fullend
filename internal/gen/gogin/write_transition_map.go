//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=states
//ff:what writes the transitions map literal for a state machine

package gogin

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// writeTransitionMap writes the transitions map literal.
func writeTransitionMap(buf *strings.Builder, transitions []statemachine.Transition) {
	buf.WriteString("var transitions = map[transitionKey]string{\n")
	for _, t := range transitions {
		buf.WriteString(fmt.Sprintf("\t{%q, %q}: %q,\n", t.From, t.Event, t.To))
	}
	buf.WriteString("}\n\n")
}
