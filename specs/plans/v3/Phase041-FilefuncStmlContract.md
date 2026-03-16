# Phase041: filefunc 전체 적용 — stml + ssac/parser + contract

## 목표

stml (parser/validator/generator) + ssac/parser + contract 5개 패키지 filefunc 준수.
17파일 → ~198파일.

## 대상 패키지

| 패키지 | 파일 | LOC | func | type | 분해 후 | 최대 파일 |
|---|---|---|---|---|---|---|
| stml/generator | 5 | 1,333 | 48 | 5 | ~53 | 516줄 |
| stml/validator | 3 | 546 | 26 | 19 | ~45 | 284줄 |
| stml/parser | 2 | 1,276 | 28 | 12 | ~40 | 705줄 |
| ssac/parser | 2 | 1,824 | 26 | 8 | ~34 | 752줄 |
| contract | 5 | 636 | 21 | 6 | ~27 | 290줄 |

합계: 17파일 (5,615 LOC, 테스트 제외) → ~198파일

## Q3 WARNING

| 패키지 | 파일 | 함수 | LOC |
|---|---|---|---|
| stml/generator | react_target.go | GeneratePage | 102 |

## 주의 사항

- stml/parser (705줄), ssac/parser (752줄) — 대형 파서. 분해 시 parse 로직 흐름 유지 주의
- stml/validator/symbol.go — Q1 depth 5, 18 type + 6 func 밀집
- contract/splice (290줄) — Q1 nesting depth 4, 분해와 동시에 flatten

## 실행 절차

1. 파일 분해 — F1/F2/F3 해소
2. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
3. Q1 해소 (contract/hash, contract/scan, contract/splice, stml/parser, ssac/parser, stml/validator/symbol depth 5 등)
4. Q3 해소 (stml/generator react_target.go GeneratePage 102줄)
5. `go build ./...` + `go test ./...`
6. `filefunc validate` 패키지 단위 확인

## 검증

1. `filefunc validate` — 대상 패키지 ERROR 0, WARNING 0
2. `go test ./...` — 전체 통과
