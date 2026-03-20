//ff:func feature=contract type=rule control=sequence
//ff:what ScanPreservedFromSource: 소스에서 보존 대상 함수를 올바르게 스캔하는지 테스트
package contract

import (
	"testing"
)

func TestScanPreservedFromSource(t *testing.T) {
	src := `package model

//fullend:preserve ssot=db/gigs.sql contract=aaa1111
func Create() {
	// custom body
}

//fullend:gen ssot=db/gigs.sql contract=bbb2222
func FindByID() {
	// generated body
}
`
	result := scanPreservedFromSource(src)
	if len(result) != 1 {
		t.Fatalf("expected 1 preserved func, got %d", len(result))
	}
	pf, ok := result["Create"]
	if !ok {
		t.Fatal("expected Create to be preserved")
	}
	if pf.Directive.Contract != "aaa1111" {
		t.Errorf("expected contract aaa1111, got %s", pf.Directive.Contract)
	}
	if !contains(pf.BodyText, "custom body") {
		t.Error("expected body to contain 'custom body'")
	}
}
