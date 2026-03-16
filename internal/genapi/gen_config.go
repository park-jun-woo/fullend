//ff:type feature=genapi type=model
//ff:what 코드 생성 설정을 보관하는 타입
package genapi

// GenConfig holds generation settings (not parsing results).
type GenConfig struct {
	ArtifactsDir string
	SpecsDir     string
	ModulePath   string
}
