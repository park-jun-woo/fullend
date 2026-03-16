//ff:type feature=ssac-gen type=model topic=model-collect
//ff:what 모델 사용 정보를 담는 구조체
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

type modelUsage struct {
	ModelName  string
	MethodName string
	Inputs     map[string]string
	Result     *parser.Result
	FuncName   string
}
