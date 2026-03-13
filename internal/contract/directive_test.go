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
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got %+v", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if got.Ownership != tt.want.Ownership || got.SSOT != tt.want.SSOT || got.Contract != tt.want.Contract {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDirectiveString(t *testing.T) {
	d := &Directive{Ownership: "gen", SSOT: "service/gig/create_gig.ssac", Contract: "a3f8c10"}
	got := d.String()
	want := "//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c10"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestDirectiveStringJS(t *testing.T) {
	d := &Directive{Ownership: "gen", SSOT: "frontend/gig_list.html", Contract: "d4e5f60"}
	got := d.StringJS()
	want := "// fullend:gen ssot=frontend/gig_list.html contract=d4e5f60"
	if got != want {
		t.Errorf("StringJS() = %q, want %q", got, want)
	}
}

func TestRoundTrip(t *testing.T) {
	original := &Directive{Ownership: "preserve", SSOT: "db/gigs.sql", Contract: "e1d9f20"}
	s := original.String()
	parsed, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(String()) error: %v", err)
	}
	if parsed.Ownership != original.Ownership || parsed.SSOT != original.SSOT || parsed.Contract != original.Contract {
		t.Errorf("roundtrip mismatch: %+v != %+v", parsed, original)
	}
}

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
