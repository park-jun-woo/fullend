//ff:func feature=rule type=generator control=iteration dimension=1
//ff:what GenerateHandler — ServiceFunc의 시퀀스를 순회하며 핸들러 body 생성
package backend

import (
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// GenerateHandler produces Go handler body code from a SSaC ServiceFunc.
// Returns the handler function body (without signature).
func GenerateHandler(fn parsessac.ServiceFunc) string {
	graph := buildSeqGraph()
	declared := map[string]string{} // var → model type
	var body strings.Builder

	for _, seq := range fn.Sequences {
		fkRef := detectFKRef(seq, declared)
		pattern := evaluateSeq(graph, seq, fn.Name, fkRef)
		body.WriteString(emitSeq(seq, pattern))
		trackDeclaration(seq, declared)
	}

	return body.String()
}
