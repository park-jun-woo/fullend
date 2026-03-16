//ff:func feature=policy type=parser control=iteration dimension=1 topic=policy-check
//ff:what 디렉토리 내 모든 .rego 파일을 파싱하여 Policy 슬라이스를 반환한다
package policy

import (
	"fmt"
	"path/filepath"
)

// ParseDir parses all .rego files in a directory.
func ParseDir(dir string) ([]*Policy, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.rego"))
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no .rego files found in %s", dir)
	}

	var policies []*Policy
	for _, path := range matches {
		p, err := ParseFile(path)
		if err != nil {
			return nil, err
		}
		policies = append(policies, p)
	}
	return policies, nil
}
