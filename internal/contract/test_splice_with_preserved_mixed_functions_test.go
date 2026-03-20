//ff:func feature=contract type=rule control=sequence
//ff:what SpliceWithPreservedMixedFunctions: 보존/생성 함수가 혼재할 때 올바르게 처리하는지 테스트
package contract

import (
	"testing"
)

func TestSpliceWithPreserved_MixedFunctions(t *testing.T) {
	newContent := `package model

//fullend:gen ssot=db/gigs.sql contract=aaa1111
func Create() {
	// gen create
}

//fullend:gen ssot=db/gigs.sql contract=bbb2222
func FindByID() {
	// gen find
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
	// Create should be preserved, FindByID should be gen.
	if !contains(result.Content, "custom create") {
		t.Error("expected Create body to be preserved")
	}
	if !contains(result.Content, "gen find") {
		t.Error("expected FindByID body to stay as generated")
	}
}
