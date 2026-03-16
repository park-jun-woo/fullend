//ff:func feature=gen-gogin type=util control=iteration
//ff:what returns true if any function uses @publish

package gogin

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

// hasPublishSequence returns true if any function uses @publish.
func hasPublishSequence(funcs []ssacparser.ServiceFunc) bool {
	for _, fn := range funcs {
		for _, seq := range fn.Sequences {
			if seq.Type == "publish" {
				return true
			}
		}
	}
	return false
}
