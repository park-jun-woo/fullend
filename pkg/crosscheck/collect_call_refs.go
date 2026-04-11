//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectCallRefs — SSaC에서 @call 참조를 수집
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func collectCallRefs(seqs []ssac.Sequence, funcName string) []callRef {
	var refs []callRef
	for _, seq := range seqs {
		if seq.Type != "call" {
			continue
		}
		if idx := strings.IndexByte(seq.Model, '.'); idx > 0 {
			refs = append(refs, callRef{
				key:     strings.ToLower(seq.Model),
				context: funcName + "/" + seq.Model,
			})
		}
	}
	return refs
}
