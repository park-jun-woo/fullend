//ff:func feature=orchestrator type=util control=sequence
//ff:what findSSaCFile locates the SSaC source file for a service function.

package orchestrator

import (
	"os"
	"path/filepath"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func findSSaCFile(sf *ssacparser.ServiceFunc, specsDir string) string {
	// Try domain structure first.
	if sf.Domain != "" {
		rel := filepath.Join("service", sf.Domain, sf.FileName)
		if _, err := os.Stat(filepath.Join(specsDir, rel)); err == nil {
			return rel
		}
	}
	// Try flat structure.
	rel := filepath.Join("service", sf.FileName)
	if _, err := os.Stat(filepath.Join(specsDir, rel)); err == nil {
		return rel
	}
	return "service/" + sf.FileName
}
