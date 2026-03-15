//ff:func feature=gen-hurl type=util
//ff:what Maps plural resource name -> first transition order index.
package hurl

import (
	"github.com/jinzhu/inflection"

	"github.com/geul-org/fullend/internal/statemachine"
)

// buildResourceFirstTransition maps plural resource name -> first transition order index.
func buildResourceFirstTransition(diagrams []*statemachine.StateDiagram, transitionOrder map[string]int) map[string]int {
	result := make(map[string]int)
	for _, d := range diagrams {
		// Pluralize diagram ID: "gig" -> "gigs".
		resource := inflection.Plural(d.ID)
		for _, t := range d.Transitions {
			if t.From == d.InitialState {
				if ord, ok := transitionOrder[t.Event]; ok {
					if existing, exists := result[resource]; !exists || ord < existing {
						result[resource] = ord
					}
				}
			}
		}
	}
	return result
}
