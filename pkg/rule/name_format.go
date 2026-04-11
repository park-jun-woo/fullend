//ff:func feature=rule type=rule control=selection
//ff:what NameFormat — 이름이 형식 규칙을 만족하는지 검증
package rule

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// NameFormat checks that a name satisfies a format rule.
// claim: string (name). Returns (true, evidence) when format violated.
func NameFormat(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*NameFormatSpec)
	c, _ := ctx.Get("claim")
	name, _ := c.(string)
	violated := false
	switch s.Pattern {
	case "uppercase-start":
		violated = len(name) > 0 && name[0] >= 'a' && name[0] <= 'z'
	case "no-dot-prefix":
		violated = strings.Contains(name, ".")
	case "dot-method":
		violated = !strings.Contains(name, ".")
	}
	if !violated {
		return false, nil
	}
	return true, &Evidence{Rule: s.Rule, Level: s.Level, Ref: name, Message: s.Message}
}
