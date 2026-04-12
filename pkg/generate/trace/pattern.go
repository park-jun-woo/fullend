//ff:type feature=rule type=model
//ff:what Pattern — Trace에서 warrant 활성화 상태를 조회하기 위한 맵
package trace

// Pattern maps warrant name → activated (true/false).
type Pattern map[string]bool
