//ff:func feature=ssac-parse type=parser control=sequence
//ff:what TestParseSubscribe: @subscribe 구독 어노테이션 파싱 후 토픽·메시지 타입 검증
package parser

import "testing"

func TestParseSubscribe(t *testing.T) {
	src := `package service

type OnOrderCompletedMessage struct {
	OrderID int64
}

// @subscribe "order.completed"
// @get Order order = Order.FindByID({ID: message.OrderID})
func OnOrderCompleted(message OnOrderCompletedMessage) {}
`
	sfs := parseTestFile(t, src)
	sf := sfs[0]
	if sf.Subscribe == nil {
		t.Fatal("expected Subscribe to be set")
	}
	assertEqual(t, "Subscribe.Topic", sf.Subscribe.Topic, "order.completed")
	assertEqual(t, "Subscribe.MessageType", sf.Subscribe.MessageType, "OnOrderCompletedMessage")
	// @subscribe는 시퀀스에 포함되지 않아야 함
	if len(sf.Sequences) != 1 {
		t.Fatalf("expected 1 sequence (subscribe filtered), got %d", len(sf.Sequences))
	}
	assertEqual(t, "seq0.Type", sf.Sequences[0].Type, SeqGet)
}
