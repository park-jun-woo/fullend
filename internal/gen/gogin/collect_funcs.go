//ff:func feature=gen-gogin type=util control=iteration
//ff:what extracts @call references without package prefix from service functions

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectFuncs extracts @call references without a package prefix.
// Package-level funcs (e.g. "auth.VerifyPassword") are called directly via import, not via Handler fields.
func collectFuncs(funcs []ssacparser.ServiceFunc) []string {
	seen := make(map[string]bool)
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			if seq.Type == "call" && seq.Model != "" {
				if !strings.Contains(seq.Model, ".") {
					seen[seq.Model] = true
				}
			}
		}
	}
	var result []string
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
