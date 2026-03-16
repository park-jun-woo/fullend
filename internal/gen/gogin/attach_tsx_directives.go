//ff:func feature=gen-gogin type=generator control=iteration
//ff:what scans pages/*.tsx files and injects // fullend:gen directive

package gogin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/geul-org/fullend/internal/contract"
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
		src, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		content := string(src)

		// Skip if directive already present.
		if strings.Contains(content, "fullend:") {
			continue
		}

		// SSOT path: STML file derives from TSX filename.
		stmlName := strings.TrimSuffix(entry.Name(), ".tsx") + ".html"
		ssotPath := "frontend/" + stmlName
		hash := contract.Hash7(content)

		d := &contract.Directive{Ownership: "gen", SSOT: ssotPath, Contract: hash}
		newContent := d.StringJS() + "\n" + content
		os.WriteFile(path, []byte(newContent), 0644)
	}
	return nil
}
