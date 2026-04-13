//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what transforms standalone service .go files by adding Server receiver and fixing imports

package gogin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

// transformServiceFiles reads each .go file in internal/service/,
// converts standalone functions to Server methods, and writes them back in place.
func transformServiceFiles(intDir string, models, funcs []string, modulePath string, doc *openapi3.T, serviceFuncs []ssacparser.ServiceFunc) error {
	serviceDir := filepath.Join(intDir, "service")
	entries, err := os.ReadDir(serviceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no service files to transform
		}
		return err
	}

	// Build filename → operationID mapping from SSaC service funcs.
	fileToOpID := buildFileToOperationID(serviceFuncs)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		path := filepath.Join(serviceDir, entry.Name())
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		opID := fileToOpID[entry.Name()]
		transformed := transformSource(string(src), models, funcs, modulePath, false, doc, opID)
		if err := os.WriteFile(path, []byte(transformed), 0644); err != nil {
			return err
		}
	}

	return nil
}
