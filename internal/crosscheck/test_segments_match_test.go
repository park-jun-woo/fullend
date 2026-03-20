//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=scenario-check
//ff:what TestSegmentsMatch: 정규화된 경로 세그먼트 쌍의 일치 여부를 테이블 테스트
package crosscheck

import "testing"

func TestSegmentsMatch(t *testing.T) {
	tests := []struct {
		a, b []string
		want bool
	}{
		{[]string{"gigs", ":param"}, []string{"gigs", ":param"}, true},
		{[]string{"gigs"}, []string{"gigs", ":param"}, false},
		{[]string{"gigs", ":param"}, []string{"users", ":param"}, false},
		{[]string{"auth", "login"}, []string{"auth", "login"}, true},
	}

	for _, tt := range tests {
		got := segmentsMatch(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("segmentsMatch(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
