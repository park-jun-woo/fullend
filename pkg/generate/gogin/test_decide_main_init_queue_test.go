//ff:func feature=gen-gogin type=test control=sequence topic=main-init
//ff:what DecideMainInit — queue backend 만으로는 Queue=false (subscriber/publish 필요)

package gogin

import "testing"

func TestDecideMainInit_QueueRequiresBackendAndWork(t *testing.T) {
	needs := DecideMainInit(MainFacts{QueueBackend: "redis"})
	if needs.Queue {
		t.Errorf("queue without subscribers/publish should be false")
	}
	if needs.NeedsContextImport {
		t.Errorf("no queue → no context via queue path")
	}
}
