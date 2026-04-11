//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkStatesDDL — @state field → DDL column 존재 검증 (X-27)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkStatesDDL(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if len(fs.DDLTables) == 0 {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkStatesDDLSeqs(g, fn.Name, fn.Sequences)...)
	}
	return errs
}
