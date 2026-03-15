//ff:func feature=gen-gogin type=util
//ff:what inserts a directive before the first func declaration

package gogin

import (
	"strings"

	"github.com/geul-org/fullend/internal/contract"
)

// injectFuncDirective inserts a directive before the first func declaration.
func injectFuncDirective(src string, d *contract.Directive) string {
	// Find "func " at the start of a line.
	idx := strings.Index(src, "\nfunc ")
	if idx >= 0 {
		return src[:idx+1] + d.String() + "\n" + src[idx+1:]
	}
	// Try at the very beginning.
	if strings.HasPrefix(src, "func ") {
		return d.String() + "\n" + src
	}
	return src
}
