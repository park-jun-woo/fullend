//ff:func feature=gen-gogin type=generator control=selection topic=states
//ff:what writes the CanTransition function based on field type

package gogin

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
)

// writeCanTransition writes the CanTransition function appropriate for the field type.
func writeCanTransition(buf *strings.Builder, d *statemachine.StateDiagram, fieldType string) {
	switch fieldType {
	case "bool":
		buf.WriteString(generateBoolCanTransition(d))
	default:
		writeStringCanTransition(buf, d)
	}
}
