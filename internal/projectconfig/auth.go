//ff:type feature=projectconfig type=model
//ff:what JWT 인증 설정 구조체
package projectconfig

type Auth struct {
	Type      string              `yaml:"type"`       // "jwt" (required when auth is present)
	SecretEnv string              `yaml:"secret_env"`
	RawClaims map[string]string   `yaml:"claims"`     // YAML original: FieldName → "claim_key" or "claim_key:go_type"
	Claims    map[string]ClaimDef `yaml:"-"`           // Parsed from RawClaims after Load()
	Roles     []string            `yaml:"roles"`       // valid role names (e.g. ["client", "freelancer"])
}
