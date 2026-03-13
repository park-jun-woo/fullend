package contract

import (
	"fmt"
	"strings"
)

// Directive represents a //fullend: ownership directive attached to generated code.
type Directive struct {
	Ownership string // "gen" or "preserve"
	SSOT      string // SSOT file relative path (e.g. "service/gig/create_gig.ssac")
	Contract  string // 7-char SHA256 hex prefix
}

// Parse extracts a Directive from a Go comment string.
// Expected format: "//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1"
func Parse(comment string) (*Directive, error) {
	s := strings.TrimSpace(comment)

	// handle both Go (//fullend:) and JS (// fullend:) comment styles
	if strings.HasPrefix(s, "//fullend:") {
		s = strings.TrimPrefix(s, "//fullend:")
	} else if strings.HasPrefix(s, "// fullend:") {
		s = strings.TrimPrefix(s, "// fullend:")
	} else {
		return nil, fmt.Errorf("not a fullend directive: %q", comment)
	}

	parts := strings.Fields(s)
	if len(parts) < 1 {
		return nil, fmt.Errorf("missing ownership in directive: %q", comment)
	}

	d := &Directive{Ownership: parts[0]}
	if d.Ownership != "gen" && d.Ownership != "preserve" {
		return nil, fmt.Errorf("invalid ownership %q: must be gen or preserve", d.Ownership)
	}

	for _, p := range parts[1:] {
		key, val, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("invalid key=value pair %q in directive", p)
		}
		switch key {
		case "ssot":
			d.SSOT = val
		case "contract":
			d.Contract = val
		default:
			return nil, fmt.Errorf("unknown directive field %q", key)
		}
	}

	if d.SSOT == "" {
		return nil, fmt.Errorf("missing ssot= in directive: %q", comment)
	}
	if d.Contract == "" {
		return nil, fmt.Errorf("missing contract= in directive: %q", comment)
	}

	return d, nil
}

// String returns the directive as a Go comment.
func (d *Directive) String() string {
	return fmt.Sprintf("//fullend:%s ssot=%s contract=%s", d.Ownership, d.SSOT, d.Contract)
}

// StringJS returns the directive as a JS comment (with space after //).
func (d *Directive) StringJS() string {
	return fmt.Sprintf("// fullend:%s ssot=%s contract=%s", d.Ownership, d.SSOT, d.Contract)
}

// IsDirective checks if a comment line is a fullend directive.
func IsDirective(comment string) bool {
	s := strings.TrimSpace(comment)
	return strings.HasPrefix(s, "//fullend:") || strings.HasPrefix(s, "// fullend:")
}
