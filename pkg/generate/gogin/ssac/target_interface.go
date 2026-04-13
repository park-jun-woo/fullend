//ff:type feature=ssac-gen type=generator topic=output
//ff:what 특정 언어의 코드 생성기가 구현해야 하는 인터페이스
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// Target은 특정 언어의 코드 생성기가 구현해야 하는 인터페이스이다.
type Target interface {
	// GenerateFunc는 하나의 서비스 함수를 타겟 언어 소스 코드로 변환한다.
	GenerateFunc(sf ssacparser.ServiceFunc, st *rule.Ground) ([]byte, error)

	// GenerateModelInterfaces는 서비스 함수에서 사용하는 모델 인터페이스를 생성한다.
	GenerateModelInterfaces(funcs []ssacparser.ServiceFunc, st *rule.Ground, outDir string) error

	// GenerateHandlerStruct는 도메인별 Handler struct를 생성한다.
	GenerateHandlerStruct(funcs []ssacparser.ServiceFunc, st *rule.Ground, outDir string) error

	// FileExtension은 생성 파일의 확장자를 반환한다. (예: ".go", ".java", ".ts")
	FileExtension() string
}

// 컴파일 타임 인터페이스 구현 확인
var _ Target = (*GoTarget)(nil)
