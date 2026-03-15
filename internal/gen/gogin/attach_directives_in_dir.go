//ff:func feature=gen-gogin type=generator
//ff:what processes all .go files in a directory to attach directives

package gogin

import (
	"os"
	"path/filepath"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// attachDirectivesInDir processes all .go files in a directory.
func attachDirectivesInDir(dir string, sfByFile map[string]ssacparser.ServiceFunc) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if err := attachDirectiveToFile(filepath.Join(dir, entry.Name()), sfByFile); err != nil {
			return err
		}
	}
	return nil
}
