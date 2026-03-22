//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=states
//ff:what generates Go state machine packages from stateDiagrams

package gogin

import (
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/statemachine"
)

// GenerateStateMachines generates Go state machine packages from stateDiagrams.
// Output: <artifactsDir>/backend/internal/states/<id>state/<id>state.go
func GenerateStateMachines(diagrams []*statemachine.StateDiagram, artifactsDir, modulePath string) error {
	if len(diagrams) == 0 {
		return nil
	}

	statesBaseDir := filepath.Join(artifactsDir, "backend", "internal", "states")

	for _, d := range diagrams {
		if err := generateSingleStateMachine(d, statesBaseDir, modulePath); err != nil {
			return err
		}
	}

	return nil
}
