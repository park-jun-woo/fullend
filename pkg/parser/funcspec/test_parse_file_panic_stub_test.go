//ff:func feature=funcspec type=parser control=sequence
//ff:what ParseFile panic stub 감지 테스트 — panic("TODO") 본문은 HasBody=false

package funcspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilePanicStub(t *testing.T) {
	dir := t.TempDir()

	src := `package billing

// @func charge
// @description 결제 처리

type ChargeRequest struct {
	Amount int
}

type ChargeResponse struct {
	TxID string
}

func Charge(req ChargeRequest) (ChargeResponse, error) {
	panic("TODO")
}
`
	path := filepath.Join(dir, "charge.go")
	os.WriteFile(path, []byte(src), 0644)

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if spec == nil {
		t.Fatal("expected non-nil spec")
	}

	if spec.HasBody {
		t.Error("HasBody = true, want false (panic stub)")
	}
}
