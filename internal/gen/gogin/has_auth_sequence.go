//ff:func feature=gen-gogin type=util control=iteration
//ff:what returns true if any function uses @auth

package gogin

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// hasAuthSequence returns true if any function uses @auth.
func hasAuthSequence(funcs []ssacparser.ServiceFunc) bool {
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			if seq.Type == "auth" {
				return true
			}
		}
	}
	return false
}
