//ff:func feature=gen-gogin type=util control=iteration dimension=2
//ff:what checks if any service function in the domain has write sequences (post/put/delete)

package gogin

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// domainNeedsDB checks if any service function in the domain has write sequences (post/put/delete).
func domainNeedsDB(serviceFuncs []ssacparser.ServiceFunc, domain string) bool {
	for _, fn := range serviceFuncs {
		if fn.Feature != domain {
			continue
		}
		for _, seq := range fn.Sequences {
			switch seq.Type {
			case ssacparser.SeqPost, ssacparser.SeqPut, ssacparser.SeqDelete:
				return true
			}
		}
	}
	return false
}
