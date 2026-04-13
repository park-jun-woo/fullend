//ff:func feature=contract type=rule control=sequence
//ff:what SpliceWithPreservedContractChange: 계약 변경 시 경고를 생성하면서 본문을 복원하는지 테스트
package contract

import (
	"testing"
)

func TestSpliceWithPreserved_ContractChange(t *testing.T) {
	newContent := `package service

//fullend:gen ssot=service/gig/create_gig.ssac contract=new1234
func CreateGig() {
	// new generated body
}
`
	preserved := map[string]*PreservedFunc{
		"CreateGig": {
			Directive: Directive{Ownership: "preserve", SSOT: "service/gig/create_gig.ssac", Contract: "old1234"},
			BodyText:  "\n\t// custom user body\n",
		},
	}

	result, err := SpliceWithPreserved(newContent, preserved, "test.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
	if result.Warnings[0].OldContract != "old1234" || result.Warnings[0].NewContract != "new1234" {
		t.Errorf("unexpected warning: %+v", result.Warnings[0])
	}
	// Body should still be restored.
	if !contains(result.Content, "custom user body") {
		t.Error("expected preserved body to be restored even with contract change")
	}
}
