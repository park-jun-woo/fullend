# Phase040: filefunc 전체 적용 — internal 소형 패키지 + gen/gogin Q3

## 목표

internal 소형 10개 패키지 filefunc 준수 + gen/gogin Q3 WARNING 4건 해소.
21파일 → ~110파일.

## 대상 패키지

| 패키지 | 파일 | LOC | func | type | 분해 후 | feature |
|---|---|---|---|---|---|---|
| cmd/fullend | 1 | 474 | 12 | 0 | ~12 | cli |
| policy | 2 | 238 | 8 | 3 | ~11 | policy |
| statemachine | 2 | 165 | 6 | 2 | ~8 | statemachine |
| funcspec | 1 | 194 | 6 | 2 | ~8 | funcspec |
| genmodel | 1 | 523 | 20 | 4 | ~24 | genmodel |
| genapi | 1 | 50 | 0 | 4 | ~4 | genapi |
| projectconfig | 1 | 186 | 2 | 13 | ~15 | projectconfig |
| gen/react | 1 | 476 | 15 | 0 | ~15 | gen-react |
| reporter | 4 | 223 | 5 | 4 | ~9 | reporter |
| scenario | 0 | 0 | 0 | 0 | 0 | scenario |

합계: 14파일 (2,529 LOC) → ~106파일

## gen/gogin Q3 WARNING (추가)

이미 filefunc 1-func/file 준수 상태. 100줄 초과 함수 4건 분해만 필요.

| 파일 | 함수 | LOC |
|---|---|---|
| generate_main.go | generateMain | 124 |
| generate_main_with_domains.go | generateMainWithDomains | 190 |
| generate_query_opts.go | generateQueryOpts | 216 |
| transform_source.go | transformSource | 105 |

hint: backtick string detected — 템플릿 문자열 별도 파일 추출 검토.

## 실행 절차

1. 파일 분해 — F1/F2/F3 해소
2. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
3. Q1 해소 (main.go, genmodel, gen/react 등 nesting depth)
4. Q3 해소 (cmd/fullend runHistory 133줄, gen/gogin 4건)
5. `go build ./...` + `go test ./...`
6. `filefunc validate` 패키지 단위 확인

## 검증

1. `filefunc validate` — 대상 패키지 ERROR 0, WARNING 0
2. `go test ./...` — 전체 통과
