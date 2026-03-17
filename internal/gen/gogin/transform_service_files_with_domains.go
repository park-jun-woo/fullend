//ff:func feature=gen-gogin type=generator control=iteration dimension=3
//ff:what transforms service files in both flat and domain subdirectories

package gogin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// transformServiceFilesWithDomains transforms service files in both flat and domain subdirectories.
func transformServiceFilesWithDomains(intDir string, serviceFuncs []ssacparser.ServiceFunc, models, funcs []string, modulePath string, doc *openapi3.T) error {
	serviceDir := filepath.Join(intDir, "service")

	// Build filename → operationID mapping from SSaC service funcs.
	fileToOpID := buildFileToOperationID(serviceFuncs)

	// Transform flat files (Domain="") directly in serviceDir.
	entries, err := os.ReadDir(serviceDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
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

	// Transform domain subdirectory files.
	domains := uniqueDomains(serviceFuncs)
	for _, domain := range domains {
		domainDir := filepath.Join(serviceDir, domain)
		entries, err := os.ReadDir(domainDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
				continue
			}
			path := filepath.Join(domainDir, entry.Name())
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			opID := fileToOpID[entry.Name()]
			transformed := transformSource(string(src), models, funcs, modulePath, true, doc, opID)
			if err := os.WriteFile(path, []byte(transformed), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
