//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkConfig — currentUser/queue 사용 시 Config 설정 필수 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkConfig(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError

	if ssacUsesCurrentUser(fs) {
		errs = append(errs, evalConfigRequired(g, "X-48", "backend.auth.claims", "currentUser used but backend.auth.claims not configured")...)
	}

	hasQueue := len(g.Pairs["SSaC.publish"]) > 0 || len(g.Pairs["SSaC.subscribe"]) > 0
	if hasQueue {
		errs = append(errs, evalConfigRequired(g, "X-56", "queue.backend", "@publish/@subscribe used but queue.backend not configured")...)
	}

	return errs
}
