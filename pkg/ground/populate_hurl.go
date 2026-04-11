//ff:func feature=rule type=loader control=sequence
//ff:what populateHurl — Hurl 엔트리에서 path, method 추출
package ground

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateHurl(g *rule.Ground, fs *fullend.Fullstack) {
	_ = g
	_ = fs
	// Hurl entries are iterated per-claim in check_hurl.go.
	// No upfront Ground population needed.
}
