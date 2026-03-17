//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=output
//ff:what inserts a file-level directive before the package declaration

package gogin

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/contract"
)

// injectFileDirective inserts a file-level directive before the package declaration.
func injectFileDirective(src string, d *contract.Directive) string {
	// Find "package " — skip any "// Code generated" comment.
	lines := strings.SplitN(src, "\n", -1)
	for i, line := range lines {
		if !strings.HasPrefix(line, "package ") {
			continue
		}
		before := strings.Join(lines[:i], "\n")
		after := strings.Join(lines[i:], "\n")
		if before != "" {
			return before + "\n" + d.String() + "\n" + after
		}
		return d.String() + "\n" + after
	}
	return d.String() + "\n" + src
}
