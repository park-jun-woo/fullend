//ff:type feature=ssac-gen type=generator topic=output
//ff:what Go 언어용 코드 생성기 구조체
package generator

import "github.com/geul-org/fullend/internal/funcspec"

// GoTarget은 Go 언어용 코드 생성기다.
type GoTarget struct {
	FuncSpecs []funcspec.FuncSpec
}
