//ff:func feature=gen-gogin type=util control=iteration dimension=2
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
