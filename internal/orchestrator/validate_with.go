//ff:func feature=orchestrator type=command control=iteration
//ff:what pre-parsed SSOT를 사용하여 검증을 실행한다
package orchestrator

import (
	"github.com/geul-org/fullend/internal/genapi"
	"github.com/geul-org/fullend/internal/reporter"
)

// ValidateWith runs validation using pre-parsed SSOTs.
func ValidateWith(root string, detected []DetectedSSOT, parsed *genapi.ParsedSSOTs, skip map[SSOTKind]bool) *reporter.Report {
	report := &reporter.Report{}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	done := make(map[SSOTKind]bool)

	// Emit steps in fixed order.
	for _, kind := range allKinds {
		if done[kind] {
			continue
		}

		// --skip takes precedence even if detected.
		if skip[kind] {
			report.Steps = append(report.Steps, reporter.StepResult{
				Name:    string(kind),
				Status:  reporter.Skip,
				Summary: "skipped (--skip)",
			})
			continue
		}

		d, ok := has[kind]
		if !ok {
			if kind == KindFunc {
				report.Steps = append(report.Steps, reporter.StepResult{
					Name:    string(kind),
					Status:  reporter.Skip,
					Summary: "no func/ directory",
				})
			} else if kind == KindStates {
				report.Steps = append(report.Steps, reporter.StepResult{
					Name:    string(kind),
					Status:  reporter.Skip,
					Summary: "no states/ directory",
				})
			} else if kind == KindPolicy {
				report.Steps = append(report.Steps, reporter.StepResult{
					Name:    string(kind),
					Status:  reporter.Skip,
					Summary: "no policy/ directory",
				})
			} else if kind == KindScenario {
				report.Steps = append(report.Steps, reporter.StepResult{
					Name:    string(kind),
					Status:  reporter.Pass,
					Summary: "no scenario tests",
					Errors:  []string{"[WARN] tests/scenario-*.hurl 파일이 없습니다 — 시나리오 테스트를 작성하세요 (--skip scenario로 억제 가능)"},
				})
			} else {
				report.Steps = append(report.Steps, reporter.StepResult{
					Name:    string(kind),
					Status:  reporter.Fail,
					Summary: "required but not found",
				})
			}
			continue
		}

		switch kind {
		case KindConfig:
			report.Steps = append(report.Steps, validateConfig(d.Path, parsed.Config))
		case KindOpenAPI:
			report.Steps = append(report.Steps, validateOpenAPI(d.Path, parsed.OpenAPIDoc))
		case KindDDL:
			report.Steps = append(report.Steps, validateDDL(root, parsed.SymbolTable))
			// Run SSaC right after DDL to reuse symbol table.
			if _, ok := has[KindSSaC]; ok {
				report.Steps = append(report.Steps, validateSSaC(root, parsed.ServiceFuncs, parsed.SymbolTable))
				done[KindSSaC] = true
			}
		case KindSSaC:
			report.Steps = append(report.Steps, validateSSaC(root, parsed.ServiceFuncs, parsed.SymbolTable))
		case KindSTML:
			report.Steps = append(report.Steps, validateSTML(root, parsed.STMLPages))
		case KindStates:
			report.Steps = append(report.Steps, validateStates(parsed.StateDiagrams, parsed.StatesErr))
		case KindPolicy:
			report.Steps = append(report.Steps, validatePolicy(parsed.Policies))
		case KindScenario:
			step, files := validateScenarioHurl(d.Path, root)
			report.Steps = append(report.Steps, step)
			parsed.HurlFiles = files
		case KindFunc:
			report.Steps = append(report.Steps, validateFunc(parsed.ProjectFuncSpecs))
		case KindModel:
			report.Steps = append(report.Steps, validateModel(d.Path))
		}
	}

	// Cross-validation step.
	report.Steps = append(report.Steps, runCrossValidate(root, parsed))

	// Contract validation step (if artifacts exist).
	report.Steps = append(report.Steps, runContractValidate(root))

	return report
}
