//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=model-collect
//ff:what extracts model names used by funcs in a specific domain

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// collectModelsForDomain extracts model names used by funcs in a specific domain.
func collectModelsForDomain(funcs []ssacparser.ServiceFunc, domain string) []string {
	seen := make(map[string]bool)
	for _, fn := range funcs {
		if fn.Domain != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			if seq.Type == "call" || seq.Model == "" {
				continue
			}
			seen[strings.SplitN(seq.Model, ".", 2)[0]] = true
		}
	}
	var result []string
	for name := range seen {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
