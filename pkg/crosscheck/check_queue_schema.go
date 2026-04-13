//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkQueueSchema — @subscribe message fields → @publish payload 스키마 (X-59)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkQueueSchema(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	pubPayloads := collectPublishPayloads(fs)
	if len(pubPayloads) == 0 {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if fn.Subscribe == nil {
			continue
		}
		errs = append(errs, checkSubscribeFields(fn, pubPayloads, g)...)
	}
	return errs
}
