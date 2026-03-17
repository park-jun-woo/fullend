//ff:func feature=cli type=command control=sequence
//ff:what SSOT 현황 요약 출력
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/orchestrator"
)

func runStatus(specsDir string) {
	detected, err := orchestrator.DetectSSOTs(specsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	lines := orchestrator.Status(specsDir, detected)
	orchestrator.PrintStatus(os.Stdout, lines)
}
