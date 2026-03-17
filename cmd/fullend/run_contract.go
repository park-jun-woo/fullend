//ff:func feature=cli type=command control=sequence
//ff:what contract 상태 검증 및 출력
package main

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/fullend/internal/contract"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func runContract(specsDir, artifactsDir string) {
	funcs, err := contract.ScanDir(artifactsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	funcs = contract.Verify(specsDir, funcs)
	reporter.PrintContract(os.Stdout, funcs)

	_, _, broken, orphan := contract.Summary(funcs)
	if broken > 0 || orphan > 0 {
		os.Exit(1)
	}
}
