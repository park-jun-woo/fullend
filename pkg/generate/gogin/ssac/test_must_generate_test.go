//ff:func feature=ssac-gen type=test-helper control=sequence
//ff:what GenerateFunc를 호출하고 실패 시 t.Fatal하는 테스트 헬퍼
package ssac

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func mustGenerate(t *testing.T, sf ssacparser.ServiceFunc, st *validator.SymbolTable) string {
	t.Helper()
	code, err := GenerateFunc(sf, st, nil)
	if err != nil {
		t.Fatalf("GenerateFunc failed: %v", err)
	}
	return string(code)
}
