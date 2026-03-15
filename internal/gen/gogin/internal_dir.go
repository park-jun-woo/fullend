//ff:func feature=gen-gogin type=util
//ff:what returns the backend/internal/ base path for artifact output

package gogin

import "path/filepath"

// internalDir returns the backend/internal/ base path.
func internalDir(artifactsDir string) string {
	return filepath.Join(artifactsDir, "backend", "internal")
}
