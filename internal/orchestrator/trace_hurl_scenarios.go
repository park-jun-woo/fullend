//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what traceHurlScenarios finds .hurl files referencing the given endpoint.

package orchestrator

import (
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func traceHurlScenarios(opID string, doc *openapi3.T, testsDir string, specsDir string) []ChainLink {
	if doc == nil || doc.Paths == nil {
		return nil
	}

	// Find the endpoint path for this operationId.
	var endpointPath string
	for path, pi := range doc.Paths.Map() {
		for _, op := range pi.Operations() {
			if op.OperationID == opID {
				endpointPath = path
				break
			}
		}
		if endpointPath != "" {
			break
		}
	}
	if endpointPath == "" {
		return nil
	}

	// Search .hurl files for the endpoint path.
	var links []ChainLink
	hurlFiles, _ := filepath.Glob(filepath.Join(testsDir, "*.hurl"))
	for _, f := range hurlFiles {
		line := grepLine(f, endpointPath)
		if line > 0 {
			relPath, _ := filepath.Rel(specsDir, f)
			links = append(links, ChainLink{
				Kind:    "Hurl",
				File:    relPath,
				Line:    line,
				Summary: "scenario: " + filepath.Base(f),
			})
		}
	}
	return links
}
