//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=http-handler
//ff:what 시퀀스에 쓰기 작업(POST/PUT/DELETE)이 있는지 확인
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func hasWriteSequence(seqs []ssacparser.Sequence) bool {
	for _, seq := range seqs {
		switch seq.Type {
		case ssacparser.SeqPost, ssacparser.SeqPut, ssacparser.SeqDelete:
			return true
		}
	}
	return false
}
