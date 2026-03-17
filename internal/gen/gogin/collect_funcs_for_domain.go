//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=model-collect
//ff:what extracts @call references without package prefix for a specific domain

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// collectFuncsForDomain extracts @call references (without package prefix) for a specific domain.
func collectFuncsForDomain(funcs []ssacparser.ServiceFunc, domain string) []string {
	seen := make(map[string]bool)
	for _, fn := range funcs {
		if fn.Domain != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			if seq.Type == "call" && seq.Model != "" && !strings.Contains(seq.Model, ".") {
				seen[seq.Model] = true
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
