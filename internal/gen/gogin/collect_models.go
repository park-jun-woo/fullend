//ff:func feature=gen-gogin type=util
//ff:what extracts unique model names from service functions

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// collectModels extracts unique model names from service functions.
func collectModels(funcs []ssacparser.ServiceFunc) []string {
	seen := make(map[string]bool)
	for _, fn := range funcs {
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
