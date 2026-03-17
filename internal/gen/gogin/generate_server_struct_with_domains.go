//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=http-handler
//ff:what creates per-domain handler.go files and central server.go

package gogin

import (
	"fmt"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// generateServerStructWithDomains creates per-domain handler.go files and central server.go.
func generateServerStructWithDomains(intDir string, serviceFuncs []ssacparser.ServiceFunc, modulePath string, doc *openapi3.T) error {
	serviceDir := filepath.Join(intDir, "service")
	domains := uniqueDomains(serviceFuncs)

	// 1. Generate per-domain handler.go.
	for _, domain := range domains {
		if err := generateDomainHandler(serviceDir, domain, serviceFuncs, modulePath); err != nil {
			return fmt.Errorf("domain %s handler: %w", domain, err)
		}
	}

	// 2. Generate central server.go.
	return generateCentralServer(serviceDir, domains, serviceFuncs, modulePath, doc)
}
