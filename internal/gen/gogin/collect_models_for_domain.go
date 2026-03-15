//ff:func feature=gen-gogin type=util
//ff:what extracts model names used by funcs in a specific domain

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectModelsForDomain extracts model names used by funcs in a specific domain.
func collectModelsForDomain(funcs []ssacparser.ServiceFunc, domain string) []string {
	seen := make(map[string]bool)
	for _, fn := range funcs {
		if fn.Domain != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			// Skip @call — package-level funcs are not models.
			if seq.Type == "call" {
				continue
			}
			if seq.Model != "" {
				parts := strings.SplitN(seq.Model, ".", 2)
				if len(parts) >= 1 {
					seen[parts[0]] = true
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
