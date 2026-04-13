//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkInputColumns — SSaC input 필드명이 DDL 컬럼에 있는지 검증 (X-13)
package crosscheck

import (
	"github.com/ettle/strcase"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkInputColumns(g *rule.Ground, funcName string, seq ssac.Sequence, table string) []CrossError {
	t, ok := g.Tables[table]
	if !ok || len(t.Columns) == 0 {
		return nil
	}
	var errs []CrossError
	for _, arg := range seq.Args {
		if arg.Source != "request" || arg.Field == "" {
			continue
		}
		snakeField := strcase.ToSnake(arg.Field)
		if _, colOk := t.Columns[snakeField]; !colOk {
			errs = append(errs, CrossError{Rule: "X-13", Context: funcName, Level: "WARNING",
				Message: "SSaC input " + arg.Field + " not found in DDL table " + table})
		}
	}
	return errs
}
