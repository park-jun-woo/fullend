//ff:func feature=policy type=parser control=iteration dimension=1
//ff:what ParsePolicyFile — 단일 .rego 파일에서 Policy 추출
package rego

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

var (
	reOwnership = regexp.MustCompile(`^#\s*@ownership\s+(\w+):\s+(\w+)\.(\w+)(?:\s+via\s+(\w+)\.(\w+))?$`)
	reAction    = regexp.MustCompile(`input\.action\s*(?:==\s*"(\w+)"|in\s*\{([^}]+)\})`)
	reResource  = regexp.MustCompile(`input\.resource\s*==\s*"(\w+)"`)
	reOwnerRef  = regexp.MustCompile(`input\.resource_owner`)
	reRoleRef   = regexp.MustCompile(`input\.(?:user|claims)\.role\s*==\s*"(\w+)"`)
	reClaimsRef = regexp.MustCompile(`input\.claims\.(\w+)`)
)

// ParsePolicyFile parses a single .rego file and extracts policy information.
func ParsePolicyFile(path string) (*Policy, []diagnostic.Diagnostic) {
	f, err := os.Open(path)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{File: path, Message: err.Error()}}
	}
	defer f.Close()

	p := &Policy{File: path}
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if om, ok := parseOwnershipLine(line); ok {
			p.Ownerships = append(p.Ownerships, om)
		}
	}

	content := strings.Join(lines, "\n")
	extractAllowRules(content, p)
	extractClaimsRefs(content, p)
	return p, nil
}
