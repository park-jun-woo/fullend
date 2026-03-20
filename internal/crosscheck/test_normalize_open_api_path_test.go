//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=scenario-check
//ff:what TestNormalizeOpenAPIPath: OpenAPI 경로를 정규화된 세그먼트로 변환하는 테이블 테스트
package crosscheck

import (
	"slices"
	"testing"
)

func TestNormalizeOpenAPIPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"/gigs/{id}", []string{"gigs", ":param"}},
		{"/gigs/{gigId}/proposals/{proposalId}", []string{"gigs", ":param", "proposals", ":param"}},
		{"/auth/login", []string{"auth", "login"}},
	}

	for _, tt := range tests {
		got := normalizeOpenAPIPath(tt.input)
		if !slices.Equal(got, tt.want) {
			t.Errorf("normalizeOpenAPIPath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
