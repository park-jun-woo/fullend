//ff:func feature=cli type=command control=sequence
//ff:what SSOT 검증 후 코드 산출
package main

import (
	"os"

	"github.com/geul-org/fullend/internal/orchestrator"
	"github.com/geul-org/fullend/internal/reporter"
)

func runGen(specsDir, artifactsDir string, skipKinds map[orchestrator.SSOTKind]bool, reset bool) {
	report, ok := orchestrator.Gen(specsDir, artifactsDir, skipKinds, reset)
	reporter.Print(os.Stdout, report)

	if !ok {
		os.Exit(1)
	}
}
