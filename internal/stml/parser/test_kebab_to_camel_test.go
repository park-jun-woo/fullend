//ff:func feature=stml-parse type=test control=iteration dimension=1
//ff:what kebab-case → camelCase 변환 검증
package parser

import "testing"

func TestKebabToCamel(t *testing.T) {
	tests := []struct{ in, want string }{
		{"project-id", "projectId"}, {"ReservationID", "ReservationID"}, {"room-id", "roomId"}, {"a-b-c", "aBC"},
	}
	for _, tt := range tests {
		got := kebabToCamel(tt.in)
		if got != tt.want { t.Errorf("kebabToCamel(%q) = %q, want %q", tt.in, got, tt.want) }
	}
}
