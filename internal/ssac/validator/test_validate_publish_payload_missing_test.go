//ff:func feature=ssac-validate type=test control=sequence
//ff:what @publish Payload 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidatePublishPayloadMissing(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "Publish", FileName: "publish.go",
		Sequences: []parser.Sequence{{Type: parser.SeqPublish, Topic: "order.completed"}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Payload 누락")
}
