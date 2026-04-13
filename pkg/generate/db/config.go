//ff:type feature=gen-gogin type=model topic=ddl
//ff:what Config — schema.sql 생성 설정 (Phase018)

package db

// Config holds schema generation configuration.
type Config struct {
	SpecsDDLDir    string // dummys/<>/specs/db/
	ArtifactsDir   string // <>/artifacts
	AutoNobodySeed bool   // opt-in: DEFAULT N FK 에 대해 INSERT seed 자동 주입
}
