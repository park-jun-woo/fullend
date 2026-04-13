//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=output
//ff:what scans pages/*.tsx files and injects // fullend:gen directive

package gogin

import (
	"os"
	"path/filepath"
	"strings"
)

// attachTSXDirectives scans pages/*.tsx files and injects // fullend:gen directive.
func attachTSXDirectives(artifactsDir string) error {
	pagesDir := filepath.Join(artifactsDir, "frontend", "src", "pages")
	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil // pages directory doesn't exist yet
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tsx") {
			continue
		}
		path := filepath.Join(pagesDir, entry.Name())
		injectTSXDirective(path, entry.Name())
	}
	return nil
}
