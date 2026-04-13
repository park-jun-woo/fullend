//ff:type feature=orchestrator type=model
//ff:what 탐지된 SSOT 파일/디렉토리 정보
package fullend

// DetectedSSOT represents a detected SSOT file or directory.
type DetectedSSOT struct {
	Kind SSOTKind
	Path string // absolute path to the relevant directory or file
}
