//ff:func feature=ssac-parse type=parser control=sequence
//ff:what @response 블록 본문 줄을 처리
package parser

// processResponseBody는 @response 블록 본문 줄을 처리한다.
func (cp *commentParser) processResponseBody(line string) {
	done, seq := handleResponseLine(line, cp.responseLines, cp.responseSuppressWarn)
	if done {
		cp.inResponse = false
		cp.sequences = append(cp.sequences, seq)
		cp.responseLines = nil
		return
	}
	cp.responseLines = append(cp.responseLines, line)
}
