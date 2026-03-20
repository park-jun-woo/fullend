//ff:func feature=ssac-validate type=test control=sequence
//ff:what @publish Topic 누락 시 ERROR 검증
package validator

import (
	"testing"
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestValidatePublishTopicMissing(t *testing.T) {
	funcs := []parser.ServiceFunc{{Name: "Publish", FileName: "publish.go",
		Sequences: []parser.Sequence{{Type: parser.SeqPublish, Inputs: map[string]string{"OrderID": "order.ID"}}},
	}}
	errs := Validate(funcs)
	assertHasError(t, errs, "Topic 누락")
}
