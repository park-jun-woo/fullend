# Mutation Test — Func 단독

### MUT-FUNC-001: @func 어노테이션 누락
- 대상: `specs/gigbridge/func/billing/charge.go`
- 변경: `// @func Charge` 주석 삭제
- 기대: WARNING — @func 어노테이션 없는 Go 파일은 funcspec으로 인식 불가
- 결과: PASS — funcspec.ParseDir에서 @func 어노테이션 기반 파싱

### MUT-FUNC-002: func body 미구현 (TODO stub)
- 대상: `specs/gigbridge/func/billing/charge.go`
- 변경: func body를 `panic("TODO")` 한 줄로 변경
- 기대: WARNING — HasBody=false로 인식, TODO 스텁 카운트 증가
- 결과: PASS — funcspec 파서가 TODO/panic 패턴 감지
