//ff:func feature=gen-gogin type=generator control=sequence topic=states
//ff:what generates CanTransition for boolean status fields

package gogin

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// generateBoolCanTransition generates CanTransition for boolean status fields.
func generateBoolCanTransition(d *statemachine.StateDiagram) string {
	// Find the mapping: which state name corresponds to true/false.
	// Convention: state name matching a common boolean pattern → true.
	// "published" → true, "unpublished" → false, etc.
	trueState, falseState := inferBoolStates(d)

	var buf strings.Builder
	buf.WriteString("// CanTransition checks if the given event is valid from the current state.\n")
	buf.WriteString("// Set DISABLE_STATE_CHECK=1 to bypass state transition checks.\n")
	buf.WriteString("func CanTransition(input Input, event string) error {\n")
	buf.WriteString("\tif os.Getenv(\"DISABLE_STATE_CHECK\") == \"1\" {\n")
	buf.WriteString("\t\treturn nil\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tcurrent := resolveState(input.Status)\n")
	buf.WriteString("\t_, ok := transitions[transitionKey{from: current, event: event}]\n")
	buf.WriteString("\tif !ok {\n")
	buf.WriteString("\t\treturn fmt.Errorf(\"cannot transition from %q via %q\", current, event)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func resolveState(v interface{}) string {\n")
	buf.WriteString("\tswitch val := v.(type) {\n")
	buf.WriteString("\tcase bool:\n")
	buf.WriteString("\t\tif val {\n")
	buf.WriteString(fmt.Sprintf("\t\t\treturn %q\n", trueState))
	buf.WriteString("\t\t}\n")
	buf.WriteString(fmt.Sprintf("\t\treturn %q\n", falseState))
	buf.WriteString("\tcase string:\n")
	buf.WriteString("\t\treturn val\n")
	buf.WriteString("\tdefault:\n")
	buf.WriteString("\t\treturn \"\"\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	return buf.String()
}
