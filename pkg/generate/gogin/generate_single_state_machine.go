//ff:func feature=gen-gogin type=generator control=sequence topic=states
//ff:what 단일 stateDiagram에서 Go 상태 머신 파일을 생성한다

package gogin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/pkg/contract"
	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
)

// generateSingleStateMachine generates a Go state machine file for one diagram.
func generateSingleStateMachine(d *statemachine.StateDiagram, statesBaseDir, modulePath string) error {
	pkgName := d.ID + "state"
	pkgDir := filepath.Join(statesBaseDir, pkgName)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("create states dir for %s: %w", d.ID, err)
	}

	src := generateStateMachineSource(d, pkgName)

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
	return nil
}
