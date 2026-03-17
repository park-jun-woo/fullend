//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 단일 ServiceFunc의 코드를 생성하고 파일로 출력
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func generateAndWrite(t Target, sf parser.ServiceFunc, outDir string, st *validator.SymbolTable) error {
	code, err := t.GenerateFunc(sf, st)
	if err != nil {
		return fmt.Errorf("%s 코드 생성 실패: %w", sf.Name, err)
	}

	ext := t.FileExtension()
	outName := strings.TrimSuffix(sf.FileName, ".ssac") + ext
	outPath := outDir
	if sf.Domain != "" {
		outPath = filepath.Join(outDir, sf.Domain)
		os.MkdirAll(outPath, 0755)
	}
	path := filepath.Join(outPath, outName)
	if err := os.WriteFile(path, code, 0644); err != nil {
		return fmt.Errorf("%s 파일 쓰기 실패: %w", path, err)
	}
	return nil
}
