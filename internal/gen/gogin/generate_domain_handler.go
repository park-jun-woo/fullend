//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=http-handler
//ff:what creates service/{domain}/handler.go with the Handler struct

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// generateDomainHandler creates service/{domain}/handler.go with the Handler struct.
func generateDomainHandler(serviceDir, domain string, serviceFuncs []ssacparser.ServiceFunc, modulePath string) error {
	domainDir := filepath.Join(serviceDir, domain)
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return err
	}

	models := collectModelsForDomain(serviceFuncs, domain)
	funcs := collectFuncsForDomain(serviceFuncs, domain)
	needsDB := domainNeedsDB(serviceFuncs, domain)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", domain))

	if needsDB {
		b.WriteString("import (\n")
		b.WriteString("\t\"database/sql\"\n\n")
		b.WriteString(fmt.Sprintf("\t\"%s/internal/model\"\n", modulePath))
		b.WriteString(")\n\n")
	} else {
		b.WriteString(fmt.Sprintf("import \"%s/internal/model\"\n\n", modulePath))
	}

	b.WriteString("// Handler handles requests for the " + domain + " domain.\n")
	b.WriteString("type Handler struct {\n")

	if needsDB {
		b.WriteString("\tDB *sql.DB\n")
	}

	for _, m := range models {
		fieldName := ucFirst(lcFirst(m) + "Model")
		b.WriteString(fmt.Sprintf("\t%s model.%sModel\n", fieldName, m))
	}

	for _, f := range funcs {
		fieldName := ucFirst(f)
		b.WriteString(fmt.Sprintf("\t%s func(args ...interface{}) (interface{}, error)\n", fieldName))
	}

	if domainNeedsJWTSecret(serviceFuncs, domain) {
		b.WriteString("\tJWTSecret string\n")
	}

	b.WriteString("}\n")

	path := filepath.Join(domainDir, "handler.go")
	return os.WriteFile(path, []byte(b.String()), 0644)
}
