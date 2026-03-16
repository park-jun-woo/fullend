//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=model-collect
//ff:what returns sorted unique non-empty domain names from service functions

package gogin

import (
	"sort"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// uniqueDomains returns sorted unique non-empty domain names.
func uniqueDomains(funcs []ssacparser.ServiceFunc) []string {
	seen := make(map[string]bool)
	for _, f := range funcs {
		if f.Domain != "" {
			seen[f.Domain] = true
		}
	}
	var result []string
	for d := range seen {
		result = append(result, d)
	}
	sort.Strings(result)
	return result
}
