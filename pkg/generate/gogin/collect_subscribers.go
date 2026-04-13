//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=model-collect
//ff:what returns service functions that have @subscribe

package gogin

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

// collectSubscribers returns service functions that have @subscribe.
func collectSubscribers(funcs []ssacparser.ServiceFunc) []ssacparser.ServiceFunc {
	var subs []ssacparser.ServiceFunc
	for _, fn := range funcs {
		if fn.Subscribe != nil {
			subs = append(subs, fn)
		}
	}
	return subs
}
