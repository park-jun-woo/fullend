# Phase042: filefunc 전체 적용 — crosscheck ✅ 완료

## 목표

crosscheck 패키지 filefunc 준수. 19파일 → ~81파일.

## 대상

| 파일 | LOC | func 수 | Q1 | Q3 |
|---|---|---|---|---|
| archived.go | 112 | 2 | depth 4 | - |
| authz.go | 51 | 1 | depth 4 | - |
| claims.go | 134 | 2 | depth 4 | - |
| crosscheck.go | 25 | 2 | - | - |
| ddl_coverage.go | 93 | 2 | depth 4 | - |
| func.go | 461 | 2 | depth 5 | 179줄 |
| func_coverage.go | 40 | 1 | depth 3 | - |
| hurl.go | 212 | 2 | depth 5 | - |
| middleware.go | 74 | 1 | depth 5 | - |
| openapi_ddl.go | 553 | 2 | depth 4 | - |
| policy.go | 138 | 1 | depth 4 | 126줄 |
| queue.go | 126 | 2 | depth 4 | 107줄 |
| roles.go | 67 | 1 | depth 3 | - |
| rules.go | 149 | 0 | - | - |
| sensitive.go | 150 | 2 | depth 4 | - |
| ssac_ddl.go | 144 | 2 | depth 3 | - |
| ssac_openapi.go | 453 | 2 | depth 5 | - |
| states.go | 176 | 2 | depth 5 | 149줄 |
| types.go | 38 | 0 | - | - |

합계: 19파일 (3,246 LOC, 테스트 제외), Q1 16건, Q3 4건

분해 후 예상: ~81파일 (76func + 5type)

## 주의 사항

- func.go (461줄), openapi_ddl.go (553줄), ssac_openapi.go (453줄) — 대형 파일. 분해 시 Check* 함수 간 공유 헬퍼 정리
- Q1이 16건으로 가장 많은 패키지 — 깊은 중첩 루프를 early-continue 또는 헬퍼 추출로 flatten
- 테스트 파일 11개 연쇄 수정

## 실행 절차

1. 파일 분해 — F1(1func/file) 해소
2. types.go F2 분해 — 1type/file
3. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
4. Q1 해소 — early-continue, 헬퍼 추출
5. Q3 해소 — CheckFuncs/CheckPolicy/CheckQueue/CheckStates 분해
6. `go build ./...` + `go test ./...`
7. `filefunc validate` 패키지 단위 확인

## 검증

1. `filefunc validate` — crosscheck ERROR 0, WARNING 0
2. `go test ./...` — 전체 통과
