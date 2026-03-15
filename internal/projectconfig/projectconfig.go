package projectconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the fullend.yaml project configuration.
type ProjectConfig struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Backend    Backend  `yaml:"backend"`
	Frontend   Frontend `yaml:"frontend"`
	Deploy     Deploy   `yaml:"deploy"`
	Session    *BuiltinBackend `yaml:"session"`
	Cache      *BuiltinBackend `yaml:"cache"`
	File       *FileBackend    `yaml:"file"`
	Queue      *QueueBackend   `yaml:"queue"`
	Authz      *AuthzConfig    `yaml:"authz"`
}

// AuthzConfig configures the authorization package.
type AuthzConfig struct {
	Package string `yaml:"package"` // custom authz package path, default: github.com/geul-org/fullend/pkg/authz
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Backend struct {
	Lang       string   `yaml:"lang"`
	Framework  string   `yaml:"framework"`
	Module     string   `yaml:"module"`
	Middleware []string `yaml:"middleware"`
	Auth       *Auth    `yaml:"auth"`
}

// ClaimDef describes a single JWT claim with its key and Go type.
type ClaimDef struct {
	Key    string // JWT claim key (e.g. "org_id")
	GoType string // Go type (e.g. "int64"), default "string"
}

type Auth struct {
	Type      string              `yaml:"type"`       // "jwt" (required when auth is present)
	SecretEnv string              `yaml:"secret_env"`
	RawClaims map[string]string   `yaml:"claims"`     // YAML original: FieldName → "claim_key" or "claim_key:go_type"
	Claims    map[string]ClaimDef `yaml:"-"`           // Parsed from RawClaims after Load()
	Roles     []string            `yaml:"roles"`       // valid role names (e.g. ["client", "freelancer"])
}

type Frontend struct {
	Lang      string `yaml:"lang"`
	Framework string `yaml:"framework"`
	Bundler   string `yaml:"bundler"`
	Name      string `yaml:"name"`
}

type Deploy struct {
	Image  string `yaml:"image"`
	Domain string `yaml:"domain"`
}

// BuiltinBackend configures session/cache backend (postgres | memory).
type BuiltinBackend struct {
	Backend string `yaml:"backend"` // "postgres" or "memory"
}

// QueueBackend configures queue backend (postgres | memory).
type QueueBackend struct {
	Backend string `yaml:"backend"` // "postgres" or "memory"
}

// FileBackend configures file storage backend (s3 | local).
type FileBackend struct {
	Backend string     `yaml:"backend"` // "s3" or "local"
	S3      *S3Config  `yaml:"s3"`
	Local   *LocalConfig `yaml:"local"`
}

type S3Config struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

type LocalConfig struct {
	Root string `yaml:"root"`
}

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
		cfg.Backend.Auth.Claims = make(map[string]ClaimDef, len(cfg.Backend.Auth.RawClaims))
		for field, raw := range cfg.Backend.Auth.RawClaims {
			parts := strings.SplitN(raw, ":", 2)
			def := ClaimDef{Key: parts[0], GoType: "string"}
			if len(parts) == 2 && parts[1] != "" {
				def.GoType = parts[1]
			}
			cfg.Backend.Auth.Claims[field] = def
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// allowedClaimTypes is the set of Go types allowed in claims definitions.
var allowedClaimTypes = map[string]bool{
	"string": true,
	"int64":  true,
	"bool":   true,
}

// jwtReservedKeys are standard JWT claim keys that must not be used as custom claim keys.
var jwtReservedKeys = map[string]bool{
	"exp": true, "iat": true, "sub": true, "iss": true,
	"aud": true, "nbf": true, "jti": true,
}

// Validate checks required fields.
func (c *ProjectConfig) Validate() error {
	if c.APIVersion != "fullend/v1" {
		return fmt.Errorf("fullend.yaml: apiVersion must be \"fullend/v1\", got %q", c.APIVersion)
	}
	if c.Kind != "Project" {
		return fmt.Errorf("fullend.yaml: kind must be \"Project\", got %q", c.Kind)
	}
	if c.Metadata.Name == "" {
		return fmt.Errorf("fullend.yaml: metadata.name is required")
	}
	if c.Backend.Module == "" {
		return fmt.Errorf("fullend.yaml: backend.module is required")
	}

	// Auth section validation.
	if c.Backend.Auth != nil {
		auth := c.Backend.Auth
		if auth.Type == "" {
			return fmt.Errorf("fullend.yaml: auth.type is required (supported: jwt)")
		}
		if auth.Type != "jwt" {
			return fmt.Errorf("fullend.yaml: auth.type %q is not supported (supported: jwt)", auth.Type)
		}
		if len(auth.Claims) == 0 {
			return fmt.Errorf("fullend.yaml: auth.claims must have at least 1 entry")
		}
		// Check each claim.
		usedKeys := make(map[string]string) // claim_key → field_name (for duplicate detection)
		for field, def := range auth.Claims {
			if !allowedClaimTypes[def.GoType] {
				return fmt.Errorf("fullend.yaml: auth.claims.%s — type %q is not allowed (allowed: string, int64, bool)", field, def.GoType)
			}
			if jwtReservedKeys[def.Key] {
				return fmt.Errorf("fullend.yaml: auth.claims.%s — claim key %q is a reserved JWT key", field, def.Key)
			}
			if prev, dup := usedKeys[def.Key]; dup {
				return fmt.Errorf("fullend.yaml: auth.claims — duplicate claim key %q (used by %s and %s)", def.Key, prev, field)
			}
			usedKeys[def.Key] = field
		}
	}

	return nil
}
