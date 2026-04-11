//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what DDL 검증 — pkg/validate/ddl + pkg/parser/ddl 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	pkgddl "github.com/park-jun-woo/fullend/pkg/validate/ddl"
)

func validateDDL(root string, st *ssacvalidator.SymbolTable) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindDDL)}

	detected, _ := fullend.DetectSSOTs(root)
	fs := fullend.ParseAll(root, detected, nil)

	tables := len(fs.DDLTables)
	cols := countPkgDDLColumns(fs)

	verrs := pkgddl.Validate(fs.DDLTables)
	for _, ve := range verrs {
		step.Errors = append(step.Errors, fmt.Sprintf("%s: %s", ve.Rule, ve.Message))
	}
	step.Errors = append(step.Errors, checkSqlcQueryDuplicates(root)...)
	step.Errors = append(step.Errors, checkDDLNullableColumns(root)...)

	if len(step.Errors) > 0 {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d tables, %d columns", tables, cols)
	return step
}
