//ff:func feature=ssac-gen type=generator control=iteration dimension=1
//ff:what 요청 파라미터 추출 코드를 버퍼에 출력
package generator

import "bytes"

func writeRequestParamsCode(buf *bytes.Buffer, requestParams []typedRequestParam) {
	for _, rp := range requestParams {
		buf.WriteString(rp.extractCode)
	}
	if len(requestParams) > 0 {
		buf.WriteString("\n")
	}
}
