//ff:type feature=gen-gogin type=model topic=main-init
//ff:what InitNeeds — main.go 초기화 블록 축별 판정 결과

package gogin

// InitNeeds carries all per-axis decisions for main.go init block emission.
type InitNeeds struct {
	Auth    bool
	Authz   bool // AUTH 과 1:1 동치 (장래 분리 대비 별도 필드)
	Queue   bool
	Session BackendNeed
	Cache   BackendNeed
	File    BackendNeed

	// NeedsContextImport derived: session=postgres ∨ cache=postgres ∨ file=s3 ∨ queue.
	NeedsContextImport bool
}
