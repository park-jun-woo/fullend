//ff:func feature=contract type=rule control=sequence
//ff:what SpliceWithPreservedNoPreserve: 보존 대상이 없을 때 원본 그대로 반환하는지 테스트
package contract

import (
	"testing"
)

func TestSpliceWithPreserved_NoPreserve(t *testing.T) {
	newContent := `package service

//fullend:gen ssot=service/gig/create_gig.ssac contract=abc1234
func CreateGig() {
	// generated body
}
`
	result, err := SpliceWithPreserved(newContent, nil, "test.go")
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != newContent {
		t.Error("expected no change when no preserves")
	}
}
