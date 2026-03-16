//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=output
//ff:what 지정된 Target으로 ServiceFunc 배열을 순회하며 코드를 생성
package generator

import (
	"fmt"
	"os"

	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

// GenerateWith는 지정된 Target으로 코드를 생성한다.
func GenerateWith(t Target, funcs []parser.ServiceFunc, outDir string, st *validator.SymbolTable) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
	}

	for _, sf := range funcs {
		if err := generateAndWrite(t, sf, outDir, st); err != nil {
			return err
		}
	}
	return nil
}
