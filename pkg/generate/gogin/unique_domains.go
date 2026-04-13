//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=model-collect
//ff:what returns sorted unique non-empty domain names from service functions

package gogin

import (
	"sort"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// uniqueDomains returns sorted unique non-empty domain names.
func uniqueDomains(funcs []ssacparser.ServiceFunc) []string {
	seen := make(map[string]bool)
	for _, f := range funcs {
		if f.Feature != "" {
			seen[f.Feature] = true
		}
	}
	var result []string
	for d := range seen {
		result = append(result, d)
	}
	sort.Strings(result)
	return result
}
