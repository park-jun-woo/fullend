//ff:func feature=contract type=rule control=iteration dimension=1
//ff:what HasFilePreserve: 파일 수준 preserve 디렉티브 존재 여부를 판별하는 테스트
package contract

import (
	"testing"
)

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
