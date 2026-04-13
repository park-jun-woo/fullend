//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=states
//ff:what writes Go const declarations for state machine states

package gogin

import (
	"fmt"
	"strings"
)

// writeStateConstants writes Go const declarations for sorted states.
func writeStateConstants(buf *strings.Builder, states []string) {
	buf.WriteString("// State constants.\nconst (\n")
	for _, s := range states {
		constName := ucFirst(s)
		buf.WriteString(fmt.Sprintf("\tState%s = %q\n", constName, s))
	}
	buf.WriteString(")\n\n")
}
