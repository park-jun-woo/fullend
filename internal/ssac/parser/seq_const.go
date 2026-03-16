//ff:type feature=ssac-parse type=model
//ff:what 시퀀스 타입 상수 및 유효성 맵
package parser

// sequence 타입 상수
const (
	SeqGet      = "get"
	SeqPost     = "post"
	SeqPut      = "put"
	SeqDelete   = "delete"
	SeqEmpty    = "empty"
	SeqExists   = "exists"
	SeqState    = "state"
	SeqAuth     = "auth"
	SeqCall     = "call"
	SeqPublish  = "publish"
	SeqResponse = "response"
)

var ValidSequenceTypes = map[string]bool{
	SeqGet:      true,
	SeqPost:     true,
	SeqPut:      true,
	SeqDelete:   true,
	SeqEmpty:    true,
	SeqExists:   true,
	SeqState:    true,
	SeqAuth:     true,
	SeqCall:     true,
	SeqPublish:  true,
	SeqResponse: true,
}
