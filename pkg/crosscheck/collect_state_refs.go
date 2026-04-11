//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectStateRefs — SSaC에서 @state 참조를 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectStateRefs(seqs []ssac.Sequence, funcName string) []stateRef {
	var refs []stateRef
	for _, seq := range seqs {
		if seq.Type == "state" {
			refs = append(refs, stateRef{diagramID: seq.DiagramID, funcName: funcName})
		}
	}
	return refs
}
