//ff:type feature=manifest type=model
//ff:what DBConfig — fullend.yaml backend.db 섹션 (DDL 파이프라인 옵션)

package manifest

// DBConfig holds optional database-related backend config.
// Currently only `auto_nobody_seed` flag for Phase018 auto seed.
type DBConfig struct {
	AutoNobodySeed bool `yaml:"auto_nobody_seed"`
}
