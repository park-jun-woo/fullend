//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what traceStates finds state diagrams referenced by @state sequences.

package orchestrator

import (
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
		links = append(links, buildStateChainLink(d.ID, specsDir, transitions))
	}
	return links
}
