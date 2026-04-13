//ff:func feature=contract type=rule control=iteration dimension=1
//ff:what IsDirective: 문자열이 fullend 디렉티브인지 판별하는 테스트
package contract

import (
	"testing"
)

func TestIsDirective(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"//fullend:gen ssot=x contract=y", true},
		{"// fullend:preserve ssot=x contract=y", true},
		{"  //fullend:gen ssot=x contract=y", true},
		{"// regular comment", false},
		{"//go:generate something", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsDirective(tt.input)
		if got != tt.want {
			t.Errorf("IsDirective(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
