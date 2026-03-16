//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what 큐 관련 import/init/subscribe 코드 블록을 생성한다

package gogin

import (
	"fmt"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func buildQueueBlocks(serviceFuncs []ssacparser.ServiceFunc, queueBackend string) (queueImport, queueInitBlock, queueSubscribeBlock string) {
	subscribers := collectSubscribers(serviceFuncs)
	needsQueue := queueBackend != "" && (len(subscribers) > 0 || hasPublishSequence(serviceFuncs))
	if !needsQueue {
		return "", "", ""
	}

	queueImport = "\n\t\"context\"\n\t\"encoding/json\"\n\t\"github.com/geul-org/fullend/pkg/queue\"\n\t\"fmt\""
	queueInitBlock = fmt.Sprintf(`
	if err := queue.Init(context.Background(), %q, conn); err != nil {
		log.Fatalf("queue init failed: %%v", err)
	}
	defer queue.Close()
`, queueBackend)

	var subLines []string
	for _, fn := range subscribers {
		if fn.Param == nil {
			continue
		}
		svcPkg := "service"
		if fn.Domain != "" {
			svcPkg = fn.Domain + "svc"
		}
		subLines = append(subLines, fmt.Sprintf(`
	queue.Subscribe(%q, func(ctx context.Context, msg []byte) error {
		var message %s.%s
		if err := json.Unmarshal(msg, &message); err != nil {
			return fmt.Errorf("unmarshal: %%w", err)
		}
		return server.%s.%s(ctx, message)
	})`, fn.Subscribe.Topic, svcPkg, fn.Param.TypeName, ucFirst(fn.Domain), fn.Name))
	}
	if len(subLines) > 0 {
		queueSubscribeBlock = strings.Join(subLines, "\n") + "\n\n\tgo queue.Start(context.Background())\n"
	}

	return queueImport, queueInitBlock, queueSubscribeBlock
}
