//ff:func feature=gen-gogin type=util
//ff:what returns true if any func in the domain uses authz or currentUser

package gogin

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// domainNeedsAuth returns true if any func in the domain uses authz or currentUser.
func domainNeedsAuth(funcs []ssacparser.ServiceFunc, domain string) bool {
	for _, fn := range funcs {
		if fn.Domain != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			if seq.Type == "auth" {
				return true
			}
			// Check if any arg references currentUser.
			for _, arg := range seq.Args {
				if arg.Source == "currentUser" {
					return true
				}
			}
		}
	}
	return false
}
