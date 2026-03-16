//ff:func feature=orchestrator type=rule control=iteration
//ff:what DDL 검증 — 테이블/컬럼 파싱 + sqlc 쿼리 중복·NOT NULL 체크
package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/reporter"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func validateDDL(root string, st *ssacvalidator.SymbolTable) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindDDL)}
	if st == nil {
		// Parse failed in ParseAll; try again for error message.
		var err error
		st, err = ssacvalidator.LoadSymbolTable(root)
		if err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("DDL/SymbolTable load error: %v", err))
			return step
		}
	}
	tables := len(st.DDLTables)
	cols := 0
	for _, t := range st.DDLTables {
		cols += len(t.Columns)
	}

	// Check sqlc query name duplicates across files.
	if dupes := checkSqlcQueryDuplicates(root); len(dupes) > 0 {
		for _, d := range dupes {
			step.Errors = append(step.Errors, d)
		}
	}

	// Check nullable columns (NOT NULL required on all columns).
	if nullables := checkDDLNullableColumns(root); len(nullables) > 0 {
		for _, n := range nullables {
			step.Errors = append(step.Errors, n)
		}
	}

	if len(step.Errors) > 0 {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d tables, %d columns", tables, cols)
	return step
}
