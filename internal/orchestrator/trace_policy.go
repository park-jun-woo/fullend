//ff:func feature=orchestrator type=util
//ff:what tracePolicy finds Rego policies referenced by @auth sequences.

package orchestrator

import (
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/policy"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func tracePolicy(sf *ssacparser.ServiceFunc, policies []*policy.Policy, specsDir string) []ChainLink {
	resources := map[string]bool{}
	actions := map[string]bool{}
	for _, seq := range sf.Sequences {
		if seq.Type != "auth" {
			continue
		}
		if seq.Resource != "" {
			resources[seq.Resource] = true
		}
		if seq.Action != "" {
			actions[seq.Action] = true
		}
	}

	if len(resources) == 0 {
		return nil
	}

	var links []ChainLink
	seen := map[string]bool{}
	for _, p := range policies {
		for _, rule := range p.Rules {
			if !resources[rule.Resource] {
				continue
			}
			relPath, _ := filepath.Rel(specsDir, p.File)
			if relPath == "" {
				relPath = p.File
			}
			if seen[relPath] {
				continue
			}
			seen[relPath] = true

			line := grepLine(p.File, rule.Resource)
			var actList []string
			for _, a := range rule.Actions {
				if actions[a] {
					actList = append(actList, a)
				}
			}
			summary := "resource: " + rule.Resource
			if len(actList) > 0 {
				summary += " [" + strings.Join(actList, ", ") + "]"
			}
			links = append(links, ChainLink{
				Kind:    "Rego",
				File:    relPath,
				Line:    line,
				Summary: summary,
			})
		}
	}
	return links
}
