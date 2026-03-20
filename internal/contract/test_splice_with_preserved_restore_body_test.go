//ff:func feature=contract type=rule control=sequence
//ff:what SpliceWithPreservedRestoreBody: 보존된 함수 본문이 올바르게 복원되는지 테스트
package contract

import (
	"testing"
)

func TestSpliceWithPreserved_RestoreBody(t *testing.T) {
	newContent := `package service

//fullend:gen ssot=service/gig/create_gig.ssac contract=abc1234
func CreateGig() {
	// new generated body
}
`
	preserved := map[string]*PreservedFunc{
		"CreateGig": {
			Directive: Directive{Ownership: "preserve", SSOT: "service/gig/create_gig.ssac", Contract: "abc1234"},
			BodyText:  "\n\t// custom user body\n",
		},
	}

	result, err := SpliceWithPreserved(newContent, preserved, "test.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(result.Warnings))
	}
	if !contains(result.Content, "custom user body") {
		t.Error("expected preserved body to be restored")
	}
	if contains(result.Content, "new generated body") {
		t.Error("expected generated body to be replaced")
	}
	if !contains(result.Content, "//fullend:preserve") {
		t.Error("expected directive to be changed to preserve")
	}
}
