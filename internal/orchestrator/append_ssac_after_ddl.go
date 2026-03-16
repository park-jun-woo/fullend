//ff:func feature=orchestrator type=util control=sequence
//ff:what DDL 검증 후 SSaC가 존재하면 즉시 SSaC 검증을 추가한다

package orchestrator

import (
	"github.com/geul-org/fullend/internal/genapi"
	"github.com/geul-org/fullend/internal/reporter"
)

// appendSSaCAfterDDL runs SSaC validation right after DDL if SSaC is detected.
func appendSSaCAfterDDL(report *reporter.Report, root string, parsed *genapi.ParsedSSOTs, has map[SSOTKind]DetectedSSOT, done map[SSOTKind]bool) {
	if _, ok := has[KindSSaC]; !ok {
		return
	}
	report.Steps = append(report.Steps, validateSSaC(root, parsed.ServiceFuncs, parsed.SymbolTable))
	done[KindSSaC] = true
}
