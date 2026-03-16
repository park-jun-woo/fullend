//ff:func feature=cli type=util control=iteration dimension=1
//ff:what --skip 플래그 파싱 및 나머지 인자 분리
package main

import (
	"fmt"
	"os"

	"github.com/geul-org/fullend/internal/orchestrator"
)

// parseSkipFlag extracts --skip flag and returns (skipKinds, remainingArgs).
func parseSkipFlag(args []string) (map[orchestrator.SSOTKind]bool, []string) {
	skipKinds := make(map[orchestrator.SSOTKind]bool)
	var remaining []string
	skip := false

	for i, a := range args {
		if skip {
			skip = false
			continue
		}
		if a != "--skip" {
			remaining = append(remaining, a)
			continue
		}
		if i+1 >= len(args) {
			fmt.Fprintln(os.Stderr, "--skip requires a comma-separated list of SSOT kinds")
			os.Exit(2)
		}
		parseSkipKinds(args[i+1], skipKinds)
		skip = true
	}

	return skipKinds, remaining
}
