//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what shallowCopyGround — Ground를 얕은 복사하여 임시 Lookup 추가 가능하게 반환
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func shallowCopyGround(g *rule.Ground) *rule.Ground {
	lookup := make(map[string]rule.StringSet, len(g.Lookup))
	for k, v := range g.Lookup {
		lookup[k] = v
	}
	schemas := make(map[string][]string, len(g.Schemas))
	for k, v := range g.Schemas {
		schemas[k] = v
	}
	return &rule.Ground{
		Lookup:  lookup,
		Types:   g.Types,
		Pairs:   g.Pairs,
		Config:  g.Config,
		Vars:    g.Vars,
		Flags:   g.Flags,
		Schemas: schemas,
	}
}
