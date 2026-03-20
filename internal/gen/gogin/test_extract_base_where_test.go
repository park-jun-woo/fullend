//ff:func feature=gen-gogin type=test control=iteration dimension=1 topic=query-opts
//ff:what extractBaseWhere: SQL WHERE절과 파라미터 수 추출 검증

package gogin

import "testing"

func TestExtractBaseWhere(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantWhere  string
		wantParams int
	}{
		{
			name:       "no where",
			sql:        "SELECT * FROM gigs ORDER BY created_at DESC;",
			wantWhere:  "",
			wantParams: 0,
		},
		{
			name:       "simple where",
			sql:        "SELECT * FROM enrollments WHERE user_id = $1 ORDER BY created_at DESC;",
			wantWhere:  "user_id = $1",
			wantParams: 1,
		},
		{
			name:       "where with two params",
			sql:        "SELECT * FROM items WHERE owner_id = $1 AND status = $2 ORDER BY id DESC;",
			wantWhere:  "owner_id = $1 AND status = $2",
			wantParams: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where, params := extractBaseWhere(tt.sql)
			if where != tt.wantWhere {
				t.Errorf("where = %q, want %q", where, tt.wantWhere)
			}
			if params != tt.wantParams {
				t.Errorf("params = %d, want %d", params, tt.wantParams)
			}
		})
	}
}
