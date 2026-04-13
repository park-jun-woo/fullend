//ff:type feature=ssac-gen type=model topic=model-collect
//ff:what 모델 사용 정보를 담는 구조체
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

type modelUsage struct {
	ModelName  string
	MethodName string
	Inputs     map[string]string
	Result     *ssacparser.Result
	FuncName   string
}
