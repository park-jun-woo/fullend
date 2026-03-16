//ff:func feature=gen-gogin type=util control=iteration dimension=1
//ff:what returns true if any service function has a non-empty Domain

package gogin

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// hasDomains returns true if any service function has a non-empty Domain.
func hasDomains(funcs []ssacparser.ServiceFunc) bool {
	for _, f := range funcs {
		if f.Domain != "" {
			return true
		}
	}
	return false
}
