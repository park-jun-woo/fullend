//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what traceStates finds state diagrams referenced by @state sequences.

package orchestrator

import (
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func traceStates(sf *ssacparser.ServiceFunc, diagrams []*statemachine.StateDiagram, specsDir string) []ChainLink {
	diagramIDs := map[string]bool{}
	transitions := map[string]string{} // diagramID -> transition name
	for _, seq := range sf.Sequences {
		if seq.Type != "state" {
			continue
		}
		diagramIDs[seq.DiagramID] = true
		transitions[seq.DiagramID] = seq.Transition
	}

	if len(diagramIDs) == 0 {
		return nil
	}

	var links []ChainLink
	for _, d := range diagrams {
		if !diagramIDs[d.ID] {
			continue
		}
		relPath := "states/" + d.ID + ".md"
		trans := transitions[d.ID]
		// Find the transition line.
		line := 0
		if trans != "" {
			line = grepLine(filepath.Join(specsDir, relPath), trans)
		}
		summary := "diagram: " + d.ID
		if trans != "" {
			summary += " -> " + trans
		}
		links = append(links, ChainLink{
			Kind:    "StateDiag",
			File:    relPath,
			Line:    line,
			Summary: summary,
		})
	}
	return links
}
