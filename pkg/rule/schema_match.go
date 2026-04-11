//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what SchemaMatch — 소스 필드들이 대상 스키마에 존재하는지 검증
package rule

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// SchemaMatch checks that source fields exist in the target schema.
// claim: []string (source field names). Returns (true, evidence) when any missing.
func SchemaMatch(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	s := specs[0].(*SchemaMatchSpec)
	g, _ := ctx.Get("ground")
	ground := g.(*Ground)
	c, _ := ctx.Get("claim")
	sourceFields, _ := c.([]string)
	targetFields, ok := ground.Schemas[s.LookupKey]
	if !ok {
		return false, nil
	}
	targetSet := make(StringSet, len(targetFields))
	for _, f := range targetFields {
		targetSet[f] = true
	}
	var missing []string
	for _, f := range sourceFields {
		if !targetSet[f] {
			missing = append(missing, f)
		}
	}
	if len(missing) == 0 {
		return false, nil
	}
	return true, &SchemaEvidence{Rule: s.Rule, Level: s.Level, Missing: missing, Message: s.Message}
}
