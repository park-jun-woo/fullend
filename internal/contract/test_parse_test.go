//ff:func feature=contract type=rule control=iteration dimension=1
//ff:what Parse: 디렉티브 문자열을 파싱하여 Directive 구조체로 변환하는 테스트
package contract

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Directive
		wantErr bool
	}{
		{
			name:  "gen directive",
			input: "//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c10",
			want:  &Directive{Ownership: "gen", SSOT: "service/gig/create_gig.ssac", Contract: "a3f8c10"},
		},
		{
			name:  "preserve directive",
			input: "//fullend:preserve ssot=db/gigs.sql contract=e1d9f20",
			want:  &Directive{Ownership: "preserve", SSOT: "db/gigs.sql", Contract: "e1d9f20"},
		},
		{
			name:  "JS style with space",
			input: "// fullend:gen ssot=frontend/gig_list.html contract=d4e5f60",
			want:  &Directive{Ownership: "gen", SSOT: "frontend/gig_list.html", Contract: "d4e5f60"},
		},
		{
			name:  "with leading whitespace",
			input: "  //fullend:gen ssot=service/gig.ssac contract=abc1234",
			want:  &Directive{Ownership: "gen", SSOT: "service/gig.ssac", Contract: "abc1234"},
		},
		{
			name:    "not a directive",
			input:   "// some regular comment",
			wantErr: true,
		},
		{
			name:    "invalid ownership",
			input:   "//fullend:modified ssot=x.ssac contract=abc1234",
			wantErr: true,
		},
		{
			name:    "missing ssot",
			input:   "//fullend:gen contract=abc1234",
			wantErr: true,
		},
		{
			name:    "missing contract",
			input:   "//fullend:gen ssot=x.ssac",
			wantErr: true,
		},
		{
			name:    "unknown field",
			input:   "//fullend:gen ssot=x.ssac contract=abc1234 foo=bar",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertParseCase(t, tt.name, tt.input, tt.want, tt.wantErr)
		})
	}
}
