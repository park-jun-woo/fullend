//ff:func feature=cli type=command control=sequence
//ff:what SSOT 검증 실행
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/orchestrator"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func runValidate(specsDir string, skipKinds map[orchestrator.SSOTKind]bool) {
	detected, err := orchestrator.DetectSSOTs(specsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	report := orchestrator.Validate(specsDir, detected, skipKinds)
	reporter.Print(os.Stdout, report)

	if report.HasFailure() {
		os.Exit(1)
	}
}
