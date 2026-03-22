//ff:func feature=orchestrator type=util control=sequence
//ff:what 단일 stateDiagram에서 ChainLink를 생성한다

package orchestrator

import "path/filepath"

// buildStateChainLink creates a ChainLink for a matched state diagram.
func buildStateChainLink(diagramID, specsDir string, transitions map[string]string) ChainLink {
	relPath := "states/" + diagramID + ".md"
	trans := transitions[diagramID]
	line := 0
	if trans != "" {
		line = grepLine(filepath.Join(specsDir, relPath), trans)
	}
	summary := "diagram: " + diagramID
	if trans != "" {
		summary += " -> " + trans
	}
	return ChainLink{
		Kind:    "StateDiag",
		File:    relPath,
		Line:    line,
		Summary: summary,
	}
}
