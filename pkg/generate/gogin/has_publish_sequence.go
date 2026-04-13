//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what returns true if any function uses @publish

package gogin

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

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
