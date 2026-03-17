//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what checks if any service function in the domain calls auth.IssueToken

package gogin

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

// domainNeedsJWTSecret checks if any service function in the domain calls auth.IssueToken.
func domainNeedsJWTSecret(serviceFuncs []ssacparser.ServiceFunc, domain string) bool {
	for _, fn := range serviceFuncs {
		if fn.Domain != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			if seq.Model == "auth.IssueToken" {
				return true
			}
		}
	}
	return false
}
