//ff:func feature=orchestrator type=command
//ff:what SSOT 검증 메인 엔트리포인트 — 개별 검증 + 교차 검증 오케스트레이션
package orchestrator

import (
	"github.com/geul-org/fullend/internal/reporter"
)

// allKinds defines the display order of SSOT kinds for validation.
var allKinds = []SSOTKind{KindConfig, KindOpenAPI, KindDDL, KindSSaC, KindModel, KindSTML, KindStates, KindPolicy, KindScenario, KindFunc}

// Validate runs individual SSOT validations on the detected sources,
// then runs cross-validation if OpenAPI + DDL + SSaC are all present.
// skipKinds specifies SSOT kinds to explicitly skip (via --skip flag).
func Validate(root string, detected []DetectedSSOT, skipKinds ...map[SSOTKind]bool) *reporter.Report {
	skip := make(map[SSOTKind]bool)
	if len(skipKinds) > 0 && skipKinds[0] != nil {
		skip = skipKinds[0]
	}

	// Parse all SSOTs once.
	parsed := ParseAll(root, detected, skip)

	return ValidateWith(root, detected, parsed, skip)
}
