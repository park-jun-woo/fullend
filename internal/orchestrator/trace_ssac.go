//ff:func feature=orchestrator type=util control=iteration
//ff:what traceSSaC locates the SSaC service function and summarizes its sequence types.

package orchestrator

import (
	"path/filepath"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func traceSSaC(sf *ssacparser.ServiceFunc, specsDir string) ChainLink {
	// Build sequence summary.
	var seqTypes []string
	seen := map[string]bool{}
	for _, seq := range sf.Sequences {
		tag := "@" + seq.Type
		if !seen[tag] {
			seqTypes = append(seqTypes, tag)
			seen[tag] = true
		}
	}

	// Find the file.
	relPath := findSSaCFile(sf, specsDir)
	line := 0
	if relPath != "" {
		line = grepLine(filepath.Join(specsDir, relPath), "func "+sf.Name)
	}

	return ChainLink{
		Kind:    "SSaC",
		File:    relPath,
		Line:    line,
		Summary: strings.Join(seqTypes, " "),
	}
}
