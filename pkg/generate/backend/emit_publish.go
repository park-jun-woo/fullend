//ff:func feature=rule type=generator control=sequence
//ff:what emitPublish — @publish 큐 발행 코드 생성
package backend

import (
	"fmt"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func emitPublish(seq parsessac.Sequence) string {
	return fmt.Sprintf("\tif err := queue.Publish(ctx, %q, %s); err != nil { return nil, err }\n",
		seq.Topic, renderFieldsAsStruct(seq))
}
