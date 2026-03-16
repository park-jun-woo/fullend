# Phase043: filefunc 전체 적용 — ssac/generator

## 목표

ssac/generator 패키지 filefunc 준수. 10파일 → ~72파일.

## 대상

| 파일 | LOC | func 수 | Q1 | Q3 |
|---|---|---|---|---|
| field_resolver.go | 58 | 0 (type 2) | depth 5 | - |
| generator.go | 132 | 3 | depth 4 | - |
| go_args.go | 185 | 4 | depth 3 | - |
| go_handler.go | 284 | 3 (method 2) | depth 4 | 145줄 |
| go_helpers.go | 487 | 6 | depth 5 | 168줄 |
| go_interface.go | 450 | 4 (type 2) | depth 5 | 108줄 |
| go_params.go | 308 | 4 | depth 5 | - |
| go_target.go | 129 | 1 (method 3) | depth 3 | - |
| go_templates.go | 183 | 0 | - | - |
| target.go | 29 | 1 | - | - |

합계: 10파일 (2,245 LOC, 테스트 제외), Q1 8건, Q3 3건

분해 후 예상: ~72파일 (62func + 10type)

## 주의 사항

- go_helpers.go (487줄) — buildTemplateData 168줄, 가장 큰 분해 대상
- go_interface.go (450줄) — deriveInterfaces 108줄 + type 2개 분리
- go_handler.go (284줄) — generateHTTPFunc 145줄, method 분리
- go_templates.go — backtick 템플릿 모음, func 없으면 F 위반 없지만 정리 대상
- 테스트 파일 4개 연쇄 수정

## 실행 절차

1. 파일 분해 — F1/F2/F3 해소
2. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
3. Q1 해소 — early-continue, 헬퍼 추출
4. Q3 해소 — generateHTTPFunc/buildTemplateData/deriveInterfaces 분해
5. `go build ./...` + `go test ./...`
6. `filefunc validate` 패키지 단위 확인

## 검증

1. `filefunc validate` — ssac/generator ERROR 0, WARNING 0
2. `go test ./...` — 전체 통과
3. `fullend validate specs/dummys/gigbridge-try02/`
4. `fullend validate specs/dummys/zenflow-try05/`
