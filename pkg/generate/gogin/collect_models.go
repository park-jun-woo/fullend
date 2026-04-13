//ff:func feature=gen-gogin type=util control=iteration dimension=2 topic=model-collect
//ff:what extracts unique model names from service functions

package gogin

import (
	"sort"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
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
