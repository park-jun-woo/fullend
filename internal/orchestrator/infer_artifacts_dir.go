//ff:func feature=orchestrator type=util
//ff:what inferArtifactsDir tries to find the artifacts directory for a specs dir.

package orchestrator

import (
	"os"
	"path/filepath"
)

// inferArtifactsDir tries to find the artifacts directory for a specs dir.
// Convention: specs/<project> -> artifacts/<project>
func inferArtifactsDir(specsDir string) string {
	base := filepath.Base(specsDir)
	candidate := filepath.Join(filepath.Dir(specsDir), "..", "artifacts", base)
	abs, err := filepath.Abs(candidate)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(abs); err == nil {
		return abs
	}
	// Also try: specsDir/../artifacts/<base>
	candidate = filepath.Join(specsDir, "..", "artifacts", base)
	abs, err = filepath.Abs(candidate)
	if err != nil {
		return ""
	}
	if _, err := os.Stat(abs); err == nil {
		return abs
	}
	return ""
}
