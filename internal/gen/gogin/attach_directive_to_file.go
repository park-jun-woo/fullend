//ff:func feature=gen-gogin type=generator
//ff:what injects a //fullend:gen directive into a single service .go file

package gogin

import (
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/contract"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// attachDirectiveToFile injects a //fullend:gen directive into a single service .go file.
func attachDirectiveToFile(path string, sfByFile map[string]ssacparser.ServiceFunc) error {
	name := filepath.Base(path)

	// Skip infrastructure files.
	if name == "handler.go" || name == "server.go" {
		return nil
	}

	sf, ok := sfByFile[name]
	if !ok {
		return nil
	}

	// Compute SSOT path.
	ssotPath := "service/" + sf.FileName
	if sf.Domain != "" {
		ssotPath = "service/" + sf.Domain + "/" + sf.FileName
	}

	d := &contract.Directive{
		Ownership: "gen",
		SSOT:      ssotPath,
		Contract:  contract.HashServiceFunc(sf),
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := injectFuncDirective(string(src), d)
	return os.WriteFile(path, []byte(content), 0644)
}
