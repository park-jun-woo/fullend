//ff:func feature=rule type=util control=iteration dimension=1
//ff:what copyGroundWithVars — Ground를 복사하고 암묵적 변수 Flags 설정
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func copyGroundWithVars(g *rule.Ground, fn parsessac.ServiceFunc) *rule.Ground {
	vars := make(rule.StringSet)
	flags := make(rule.StringSet, len(g.Flags))
	for k, v := range g.Flags {
		flags[k] = v
	}
	// implicit variables
	flags["implicit.request"] = true
	flags["implicit.currentUser"] = true
	flags["implicit.query"] = true
	if fn.Subscribe != nil {
		flags["implicit.message"] = true
	}
	return &rule.Ground{
		Lookup: g.Lookup, Types: g.Types, Pairs: g.Pairs,
		Config: g.Config, Vars: vars, Flags: flags, Schemas: g.Schemas,
	}
}
