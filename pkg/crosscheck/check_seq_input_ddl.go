//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkSeqInputDDL — 시퀀스 input의 DDL 컬럼 존재 + key case 검증 (X-13, X-14)
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkSeqInputDDL(g *rule.Ground, funcName string, seq ssac.Sequence) []CrossError {
	var errs []CrossError
	model := extractModelName(seq.Model)
	if model != "" {
		table := strings.ToLower(inflection.Plural(model))
		errs = append(errs, checkInputColumns(g, funcName, seq, table)...)
	}
	errs = append(errs, checkInputKeyCase(funcName, seq.Args)...)
	return errs
}
