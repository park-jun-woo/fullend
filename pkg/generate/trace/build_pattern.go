//ff:func feature=rule type=util control=iteration dimension=1
//ff:what BuildPattern — toulmin EvalResult에서 warrant 활성화 맵 구축
package trace

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// BuildPattern extracts activated warrants from evaluation results.
func BuildPattern(results []toulmin.EvalResult) Pattern {
	p := make(Pattern)
	for _, r := range results {
		mergeTrace(p, r.Trace)
	}
	return p
}
