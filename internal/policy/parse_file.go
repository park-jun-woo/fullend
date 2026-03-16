//ff:func feature=policy type=parser control=iteration dimension=1 topic=policy-check
//ff:what .rego 파일 하나를 파싱하여 Policy 구조체를 반환한다
package policy

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reOwnership = regexp.MustCompile(
		`^#\s*@ownership\s+(\w+):\s+(\w+)\.(\w+)(?:\s+via\s+(\w+)\.(\w+))?$`,
	)
	reAction   = regexp.MustCompile(`input\.action\s*(?:==\s*"(\w+)"|in\s*\{([^}]+)\})`)
	reResource = regexp.MustCompile(`input\.resource\s*==\s*"(\w+)"`)
	reOwnerRef = regexp.MustCompile(`input\.resource_owner`)
	reRoleRef   = regexp.MustCompile(`input\.(?:user|claims)\.role\s*==\s*"(\w+)"`)
	reClaimsRef = regexp.MustCompile(`input\.claims\.(\w+)`)
)

// ParseFile parses a single .rego file and extracts policy information.
func ParseFile(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	p := &Policy{File: path}
	scanner := bufio.NewScanner(f)

	// First pass: extract @ownership annotations.
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if om, ok := parseOwnershipLine(line); ok {
			p.Ownerships = append(p.Ownerships, om)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	// Second pass: extract allow rules by scanning allow blocks.
	content := strings.Join(lines, "\n")
	extractAllowRules(content, p)

	// Extract all input.claims.xxx references (deduplicated).
	seen := make(map[string]bool)
	for _, m := range reClaimsRef.FindAllStringSubmatch(content, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			p.ClaimsRefs = append(p.ClaimsRefs, m[1])
		}
	}

	return p, nil
}
