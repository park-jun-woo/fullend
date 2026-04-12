//ff:func feature=rule type=util control=iteration dimension=1
//ff:what mergeTrace — TraceEntry 목록에서 warrant 활성화를 Pattern에 병합
package trace

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func mergeTrace(p Pattern, entries []toulmin.TraceEntry) {
	for _, t := range entries {
		if t.Role == "rule" {
			p[t.Name] = t.Activated
		}
	}
}
