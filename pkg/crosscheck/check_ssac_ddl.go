//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkSSaCDDL — SSaC @result 타입 ↔ DDL, input ↔ DDL 교차 검증 (X-11~X-13)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkSSaCDDL(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		for _, seq := range fn.Sequences {
			errs = append(errs, checkSeqDDL(g, fn.Name, seq)...)
		}
	}
	return errs
}
