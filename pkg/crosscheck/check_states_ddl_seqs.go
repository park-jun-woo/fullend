//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkStatesDDLSeqs — 단일 함수의 @state 시퀀스별 DDL 컬럼 존재 검증 (X-27)
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkStatesDDLSeqs(g *rule.Ground, funcName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "state" || seq.DiagramID == "" {
			continue
		}
		table := strings.ToLower(inflection.Plural(seq.DiagramID))
		for field := range seq.Inputs {
			errs = append(errs, evalColumnRef(g, table, field, "X-27", funcName)...)
		}
	}
	return errs
}
