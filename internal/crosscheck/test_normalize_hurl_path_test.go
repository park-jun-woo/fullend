//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=scenario-check
//ff:what TestNormalizeHurlPath: Hurl 경로를 정규화된 세그먼트로 변환하는 테이블 테스트
package crosscheck

import (
	"slices"
	"testing"
)

func TestNormalizeHurlPath(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"/gigs/{{gig_id}}", []string{"gigs", ":param"}},
		{"/gigs", []string{"gigs"}},
		{"/gigs/{{gig_id}}/proposals", []string{"gigs", ":param", "proposals"}},
		{"/auth/login", []string{"auth", "login"}},
		{"/gigs?status=open", []string{"gigs"}},
	}

	for _, tt := range tests {
		got := normalizeHurlPath(tt.input)
		if !slices.Equal(got, tt.want) {
			t.Errorf("normalizeHurlPath(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
