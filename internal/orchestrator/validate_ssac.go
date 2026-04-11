//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what SSaC 서비스 함수 검증 — pkg/validate/ssac toulmin 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
	pkgssac "github.com/park-jun-woo/fullend/pkg/validate/ssac"
)

func validateSSaC(root string, funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindSSaC)}
	if funcs == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "SSaC parse failed")
		return step
	}

	// Re-parse via pkg/parser for pkg types
	detected, _ := fullend.DetectSSOTs(root)
	fs := fullend.ParseAll(root, detected, nil)
	ground := &rule.Ground{
		Lookup:  make(map[string]rule.StringSet),
		Types:   make(map[string]string),
		Pairs:   make(map[string]rule.StringSet),
		Config:  make(map[string]bool),
		Vars:    make(rule.StringSet),
		Flags:   make(rule.StringSet),
		Schemas: make(map[string][]string),
	}
	// Populate Go reserved words for ForbiddenRef
	ground.Lookup["go.reserved"] = rule.StringSet{
		"break": true, "case": true, "chan": true, "const": true, "continue": true,
		"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
		"func": true, "go": true, "goto": true, "if": true, "import": true,
		"interface": true, "map": true, "package": true, "range": true, "return": true,
		"select": true, "struct": true, "switch": true, "type": true, "var": true,
	}

	verrs := pkgssac.Validate(fs.ServiceFuncs, ground)
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
