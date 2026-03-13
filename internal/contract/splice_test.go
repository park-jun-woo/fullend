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

func TestHasFilePreserve(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "file-level preserve",
			src:  "//fullend:preserve ssot=states/gig.md contract=abc1234\npackage gigstate\n",
			want: true,
		},
		{
			name: "file-level gen",
			src:  "//fullend:gen ssot=states/gig.md contract=abc1234\npackage gigstate\n",
			want: false,
		},
		{
			name: "no directive",
			src:  "package service\n\nfunc Foo() {}\n",
			want: false,
		},
		{
			name: "code gen comment then preserve",
			src:  "// Code generated — do not edit.\n//fullend:preserve ssot=states/gig.md contract=abc1234\npackage gigstate\n",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasFilePreserve(tt.src); got != tt.want {
				t.Errorf("hasFilePreserve() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func contains(s, sub string) bool {
	return len(s) >= len(sub) && containsStr(s, sub)
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
