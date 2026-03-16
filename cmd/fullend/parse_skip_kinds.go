//ff:func feature=cli type=util control=iteration dimension=1
//ff:what 쉼표 구분 SSOT 종류 문자열을 파싱
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/geul-org/fullend/internal/orchestrator"
)

// parseSkipKinds parses a comma-separated string of SSOT kinds into a map.
func parseSkipKinds(csv string, skipKinds map[orchestrator.SSOTKind]bool) {
	for _, s := range strings.Split(csv, ",") {
		s = strings.TrimSpace(s)
		kind, ok := orchestrator.KindFromString(s)
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown SSOT kind: %q\nvalid kinds: openapi, ddl, ssac, model, stml, states, policy, scenario, func\n", s)
			os.Exit(2)
		}
		skipKinds[kind] = true
	}
}
