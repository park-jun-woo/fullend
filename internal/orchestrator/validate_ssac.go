//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what SSaC 서비스 함수 검증 — pkg/validate/ssac + pkg/ground 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
	"github.com/park-jun-woo/fullend/pkg/ground"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	pkgssac "github.com/park-jun-woo/fullend/pkg/validate/ssac"
)

func validateSSaC(root string, funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindSSaC)}
	if funcs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "SSaC parse failed")
		return step
	}

	detected, _ := fullend.DetectSSOTs(root)
	fs := fullend.ParseAll(root, detected, nil)
	g := ground.Build(fs)

	verrs := pkgssac.Validate(fs.ServiceFuncs, g)
	hasError := false
	for _, ve := range verrs {
		prefix := ""
		if ve.Level == "WARNING" {
			prefix = "[WARN] "
		} else {
			hasError = true
		}
		step.Errors = append(step.Errors, fmt.Sprintf("%s%s:%s seq[%d] — %s",
			prefix, ve.File, ve.Func, ve.SeqIdx, ve.Message))
	}
	if hasError {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d service functions", len(funcs))
	return step
}
