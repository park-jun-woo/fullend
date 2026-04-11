//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkSeqDDL — 개별 시퀀스의 result/input DDL 검증 (X-11~X-14)
package crosscheck

import (
	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkSeqDDL(g *rule.Ground, funcName string, seq ssac.Sequence) []CrossError {
	var errs []CrossError

	// X-11: result type plural form WARNING
	if seq.Result != nil && seq.Result.Type != "" {
		plural := inflection.Plural(seq.Result.Type)
		if seq.Result.Type == plural && seq.Result.Wrapper == "" {
			errs = append(errs, CrossError{Rule: "X-11", Context: funcName, Level: "WARNING",
				Message: "@result type " + seq.Result.Type + " appears to be plural"})
		}
	}

	// X-12: result type has no DDL table WARNING
	if seq.Result != nil && seq.Result.Type != "" {
		errs = append(errs, checkResultDDLTable(g, funcName, seq.Result.Type)...)
	}

	// X-13 + X-14: input column/case checks
	if seq.Type == "get" || seq.Type == "post" || seq.Type == "put" || seq.Type == "delete" {
		errs = append(errs, checkSeqInputDDL(g, funcName, seq)...)
	}

	return errs
}
