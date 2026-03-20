//ff:func feature=contract type=rule control=sequence
//ff:what SpliceWithPreservedNewFunction: 새 함수 추가 시 보존 함수와 공존하는지 테스트
package contract

import (
	"testing"
)

func TestSpliceWithPreserved_NewFunction(t *testing.T) {
	newContent := `package model

//fullend:gen ssot=db/gigs.sql contract=aaa1111
func Create() {
	// gen create
}

//fullend:gen ssot=db/gigs.sql contract=ccc3333
func Delete() {
	// gen delete (new function)
}
`
	preserved := map[string]*PreservedFunc{
		"Create": {
			Directive: Directive{Ownership: "preserve", SSOT: "db/gigs.sql", Contract: "aaa1111"},
			BodyText:  "\n\t// custom create\n",
		},
	}

	result, err := SpliceWithPreserved(newContent, preserved, "test.go")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(result.Content, "custom create") {
		t.Error("expected Create body to be preserved")
	}
	if !contains(result.Content, "gen delete") {
		t.Error("expected new Delete function to be added")
	}
}
