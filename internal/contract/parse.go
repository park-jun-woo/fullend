//ff:func feature=contract type=util control=sequence
//ff:what 코멘트 문자열에서 fullend 디렉티브를 파싱한다
package contract

import (
	"fmt"
	"strings"
)

// Parse extracts a Directive from a Go comment string.
// Expected format: "//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1"
func Parse(comment string) (*Directive, error) {
	s, err := stripDirectivePrefix(comment)
	if err != nil {
		return nil, err
	}

	parts := strings.Fields(s)
	if len(parts) < 1 {
		return nil, fmt.Errorf("missing ownership in directive: %q", comment)
	}

	d := &Directive{Ownership: parts[0]}
	if d.Ownership != "gen" && d.Ownership != "preserve" {
		return nil, fmt.Errorf("invalid ownership %q: must be gen or preserve", d.Ownership)
	}

	if err := parseDirectiveFields(d, parts[1:]); err != nil {
		return nil, err
	}

	if d.SSOT == "" {
		return nil, fmt.Errorf("missing ssot= in directive: %q", comment)
	}
	if d.Contract == "" {
		return nil, fmt.Errorf("missing contract= in directive: %q", comment)
	}

	return d, nil
}
