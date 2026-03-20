//ff:func feature=contract type=rule control=sequence
//ff:what DirectiveStringJS: StringJS() 메서드가 JS 스타일 디렉티브 문자열을 생성하는지 테스트
package contract

import (
	"testing"
)

func TestDirectiveStringJS(t *testing.T) {
	d := &Directive{Ownership: "gen", SSOT: "frontend/gig_list.html", Contract: "d4e5f60"}
	got := d.StringJS()
	want := "// fullend:gen ssot=frontend/gig_list.html contract=d4e5f60"
	if got != want {
		t.Errorf("StringJS() = %q, want %q", got, want)
	}
}
