//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 심볼 테이블과 SSaC spec을 교차하여 models_gen.go를 생성
package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

// GenerateModelInterfaces는 심볼 테이블과 SSaC spec을 교차하여 Model interface를 생성한다.
func (g *GoTarget) GenerateModelInterfaces(funcs []parser.ServiceFunc, st *validator.SymbolTable, outDir string) error {
	modelDir := filepath.Join(outDir, "model")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("model 디렉토리 생성 실패: %w", err)
	}

	usages := collectModelUsages(funcs)
	interfaces := deriveInterfaces(usages, st)
	if len(interfaces) == 0 {
		return nil
	}

	code := renderInterfaces(interfaces)
	formatted, err := format.Source(code)
	if err != nil {
		return fmt.Errorf("models_gen.go gofmt 실패: %w\n--- raw ---\n%s", err, string(code))
	}

	path := filepath.Join(modelDir, "models_gen.go")
	return os.WriteFile(path, formatted, 0644)
}
