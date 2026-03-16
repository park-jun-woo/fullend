//ff:func feature=gen-gogin type=generator control=iteration
//ff:what generates Go state machine packages from stateDiagrams

package gogin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/statemachine"
)

// GenerateStateMachines generates Go state machine packages from stateDiagrams.
// Output: <artifactsDir>/backend/internal/states/<id>state/<id>state.go
func GenerateStateMachines(diagrams []*statemachine.StateDiagram, artifactsDir, modulePath string) error {
	if len(diagrams) == 0 {
		return nil
	}

	statesBaseDir := filepath.Join(artifactsDir, "backend", "internal", "states")

	for _, d := range diagrams {
		pkgName := d.ID + "state"
		pkgDir := filepath.Join(statesBaseDir, pkgName)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			return fmt.Errorf("create states dir for %s: %w", d.ID, err)
		}

		src := generateStateMachineSource(d, pkgName)

		// Inject file-level //fullend:gen directive.
		dir := &contract.Directive{
			Ownership: "gen",
			SSOT:      "states/" + d.ID + ".md",
			Contract:  contract.HashStateDiagram(d),
		}
		src = injectFileDirective(src, dir)

		outPath := filepath.Join(pkgDir, pkgName+".go")
		if err := os.WriteFile(outPath, []byte(src), 0644); err != nil {
			return fmt.Errorf("write state machine %s: %w", d.ID, err)
		}
	}

	return nil
}
