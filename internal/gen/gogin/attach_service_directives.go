//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=output
//ff:what scans service .go files and injects //fullend:gen directives

package gogin

import (
	"os"
	"path/filepath"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// attachServiceDirectives scans service .go files and injects //fullend:gen directives.
func attachServiceDirectives(intDir string, serviceFuncs []ssacparser.ServiceFunc) error {
	// Build map: .go filename → ServiceFunc
	sfByFile := make(map[string]ssacparser.ServiceFunc)
	for _, sf := range serviceFuncs {
		goFile := strings.TrimSuffix(sf.FileName, ".ssac") + ".go"
		sfByFile[goFile] = sf
	}

	serviceDir := filepath.Join(intDir, "service")

	// Process flat files and domain subdirectories.
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		return nil // no service dir yet
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := attachDirectivesInDir(filepath.Join(serviceDir, entry.Name()), sfByFile); err != nil {
			return err
		}
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if err := attachDirectiveToFile(filepath.Join(serviceDir, entry.Name()), sfByFile); err != nil {
			return err
		}
	}

	return nil
}
