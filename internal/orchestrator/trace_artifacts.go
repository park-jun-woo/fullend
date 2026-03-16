//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what traceArtifacts finds generated code artifacts connected to the operationId.

package orchestrator

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/geul-org/fullend/internal/contract"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// traceArtifacts finds generated code artifacts connected to the operationId.
func traceArtifacts(artifactsDir, operationID string, sf *ssacparser.ServiceFunc) []ChainLink {
	var links []ChainLink

	funcs, err := contract.ScanDir(artifactsDir)
	if err != nil {
		return links
	}

	// Build SSOT path for matching.
	ssotPath := "service/" + sf.FileName
	if sf.Domain != "" {
		ssotPath = "service/" + sf.Domain + "/" + sf.FileName
	}

	for _, f := range funcs {
		if f.Directive.SSOT != ssotPath {
			continue
		}
		kind := "Handler"
		if strings.Contains(f.File, "/model/") {
			kind = "Model"
		} else if strings.Contains(f.File, "/authz/") {
			kind = "Authz"
		} else if strings.Contains(f.File, "/states/") {
			kind = "States"
		}
		links = append(links, ChainLink{
			Kind:      kind,
			File:      f.File,
			Summary:   f.Function,
			Ownership: f.Status,
		})
	}

	// Also trace model methods for tables used by this operation.
	for _, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == "call" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) != 2 {
			continue
		}
		modelName := parts[0]
		methodName := parts[1]

		tableName := inflection.Plural(strings.ToLower(modelName))
		for _, f := range funcs {
			if f.Function != methodName || !strings.Contains(f.File, "/model/") || !strings.Contains(f.Directive.SSOT, tableName) {
				continue
			}
			links = append(links, ChainLink{
				Kind:      "Model",
				File:      f.File,
				Summary:   modelName + "." + methodName,
				Ownership: f.Status,
			})
		}
	}

	// Deduplicate.
	seen := make(map[string]bool)
	var unique []ChainLink
	for _, l := range links {
		key := l.Kind + "|" + l.File + "|" + l.Summary
		if !seen[key] {
			seen[key] = true
			unique = append(unique, l)
		}
	}

	return unique
}
