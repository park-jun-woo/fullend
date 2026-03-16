//ff:func feature=projectconfig type=loader control=sequence
//ff:what fullend.yaml 파일을 읽어 파싱하고 검증한 뒤 ProjectConfig를 반환한다
package projectconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load reads and parses fullend.yaml from the given specs directory root.
func Load(specsDir string) (*ProjectConfig, error) {
	path := filepath.Join(specsDir, "fullend.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("fullend.yaml not found: %w", err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("fullend.yaml parse error: %w", err)
	}

	// Post-process: convert RawClaims → Claims (ClaimDef).
	if cfg.Backend.Auth != nil && len(cfg.Backend.Auth.RawClaims) > 0 {
		cfg.Backend.Auth.Claims = parseRawClaims(cfg.Backend.Auth.RawClaims)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
