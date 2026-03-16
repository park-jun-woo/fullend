//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what traceSTML finds STML frontend files referencing the given operationId.

package orchestrator

import (
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func traceSTML(doc *openapi3.T, opID string, stmlDir string, specsDir string) []ChainLink {
	// STML references operationId via data-fetch="OpID" or data-action="OpID".
	// Search for the operationId directly in STML files.
	var links []ChainLink
	matches, _ := filepath.Glob(filepath.Join(stmlDir, "*.html"))
	for _, m := range matches {
		line := grepLine(m, opID)
		if line > 0 {
			relPath, _ := filepath.Rel(specsDir, m)
			// Determine if it's fetch or action from the matched line.
			attr := stmlMatchAttr(m, opID)
			links = append(links, ChainLink{
				Kind:    "STML",
				File:    relPath,
				Line:    line,
				Summary: attr + "=\"" + opID + "\"",
			})
		}
	}
	return links
}
