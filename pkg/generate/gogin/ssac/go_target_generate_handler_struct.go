//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=output
//ff:what 도메인별 Handler struct를 생성하여 파일로 출력
package ssac

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/ettle/strcase"
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// GenerateHandlerStruct는 도메인별 Handler struct를 생성한다.
func (g *GoTarget) GenerateHandlerStruct(funcs []ssacparser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	domainModels := collectDomainModels(funcs)
	for domain, models := range domainModels {
		if err := writeHandlerFile(domain, models, outDir); err != nil {
			return err
		}
	}
	return nil
}

func writeHandlerFile(domain string, models []string, outDir string) error {
	if len(models) == 0 {
		return nil
	}

	var buf bytes.Buffer
	pkgName := "service"
	if domain != "" {
		pkgName = domain
	}
	buf.WriteString("package " + pkgName + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"database/sql\"\n\n")
	buf.WriteString("\t\"model\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("type Handler struct {\n")
	buf.WriteString("\tDB *sql.DB\n")
	for _, m := range models {
		pascalName := strcase.ToGoPascal(m)
		fmt.Fprintf(&buf, "\t%sModel model.%sModel\n", pascalName, pascalName)
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("handler.go gofmt 실패: %w\n--- raw ---\n%s", err, buf.String())
	}

	outPath := outDir
	if domain != "" {
		outPath = filepath.Join(outDir, domain)
		os.MkdirAll(outPath, 0755)
	}
	path := filepath.Join(outPath, "handler.go")
	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return fmt.Errorf("handler.go 파일 쓰기 실패: %w", err)
	}
	return nil
}
