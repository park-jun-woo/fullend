//ff:func feature=orchestrator type=rule control=sequence
//ff:what SSaC 서비스 함수 검증 — 심볼 테이블 기반 시퀀스 유효성 검사
package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/reporter"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func validateSSaC(root string, funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindSSaC)}
	if funcs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "SSaC parse failed")
		return step
	}

	if st == nil {
		var stErr error
		st, stErr = ssacvalidator.LoadSymbolTable(root)
		if stErr != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("SSaC symbol table load error: %v", stErr))
			return step
		}
	}

	verrs := ssacvalidator.ValidateWithSymbols(funcs, st)
	if len(verrs) > 0 {
		hasError := false
		for _, ve := range verrs {
			prefix := ""
			if ve.Level == "WARNING" {
				prefix = "[WARN] "
			} else {
				hasError = true
			}
			step.Errors = append(step.Errors, fmt.Sprintf("%s%s:%s seq[%d] %s — %s",
				prefix, ve.FileName, ve.FuncName, ve.SeqIndex, ve.Tag, ve.Message))
		}
		if hasError {
			step.Status = reporter.Fail
		} else {
			step.Status = reporter.Pass
		}
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d service functions", len(funcs))
	return step
}
