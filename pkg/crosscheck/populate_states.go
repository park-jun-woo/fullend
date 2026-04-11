//ff:func feature=crosscheck type=loader control=iteration dimension=1
//ff:what populateStates — StateDiagram에서 diagram ID, event 추출
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func populateStates(g *rule.Ground, fs *fullend.Fullstack) {
	diagrams := make(rule.StringSet)
	for _, sd := range fs.StateDiagrams {
		diagrams[sd.ID] = true
		events := make(rule.StringSet)
		for _, tr := range sd.Transitions {
			events[tr.Event] = true
		}
		g.Lookup["States.event."+sd.ID] = events
	}
	g.Lookup["States.diagram"] = diagrams
}
