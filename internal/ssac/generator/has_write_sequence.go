//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=http-handler
//ff:what 시퀀스에 쓰기 작업(POST/PUT/DELETE)이 있는지 확인
package generator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

func hasWriteSequence(seqs []parser.Sequence) bool {
	for _, seq := range seqs {
		switch seq.Type {
		case parser.SeqPost, parser.SeqPut, parser.SeqDelete:
			return true
		}
	}
	return false
}
