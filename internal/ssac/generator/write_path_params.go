//ff:func feature=ssac-gen type=generator control=iteration dimension=1 topic=path-params
//ff:what 경로 파라미터 추출 코드를 버퍼에 출력
package generator

import (
	"bytes"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

func writePathParams(buf *bytes.Buffer, pathParams []validator.PathParam) {
	for _, pp := range pathParams {
		buf.WriteString(generatePathParamCode(pp))
	}
	if len(pathParams) > 0 {
		buf.WriteString("\n")
	}
}
