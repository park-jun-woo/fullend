//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkStatesGuard — 상태 전이에 참여하는 SSaC 함수에 @state 있는지 WARNING (X-26)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkStatesGuard(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.StateDiagrams) == 0 {
		return nil
	}
	transitionFuncs := make(rule.StringSet)
	for _, sd := range fs.StateDiagrams {
		for _, tr := range sd.Transitions {
			transitionFuncs[tr.Event] = true
		}
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if !transitionFuncs[fn.Name] {
			continue
		}
		if !hasStateSeq(fn.Sequences) {
			errs = append(errs, CrossError{Rule: "X-26", Context: fn.Name, Level: "WARNING",
				Message: "function participates in state transition but has no @state sequence"})
		}
	}
	return errs
}
