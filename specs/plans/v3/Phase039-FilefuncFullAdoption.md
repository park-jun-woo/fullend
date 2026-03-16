# Phase039: filefunc 전체 적용 — pkg/* 유틸리티 패키지

## 목표

pkg/ 하위 13개 패키지 filefunc 준수. 26파일 → ~107파일.

## 현황 (pkg/ 한정)

| 위반 유형 | 건수 |
|---|---|
| A1 | 26 |
| A3 | 26 |
| F1 | 4 |
| F2 | 16 |
| F3 | 3 |
| Q1 | 1 |

## 대상 패키지

| 패키지 | 파일 | LOC | func | type | 분해 후 |
|---|---|---|---|---|---|
| pkg/auth | 3 | 61 | 3 | 6 | ~9 |
| pkg/authz | 1 | 135 | 3 | 3 | ~6 |
| pkg/cache | 1 | 111 | 8 | 4 | ~12 |
| pkg/config | 1 | 18 | 2 | 0 | ~2 |
| pkg/crypto | 4 | 147 | 4 | 8 | ~12 |
| pkg/file | 1 | 93 | 8 | 3 | ~11 |
| pkg/image | 2 | 62 | 2 | 4 | ~6 |
| pkg/mail | 2 | 75 | 2 | 4 | ~6 |
| pkg/pagination | 2 | 15 | 0 | 2 | ~2 |
| pkg/queue | 1 | 261 | 8 | 2 | ~10 |
| pkg/session | 1 | 111 | 8 | 4 | ~12 |
| pkg/storage | 4 | 148 | 4 | 6 | ~10 |
| pkg/text | 3 | 63 | 3 | 6 | ~9 |

합계: 26파일 (1,300 LOC) → ~107파일

## 실행 절차

1. 파일 분해 — F1(1func/file), F2(1type/file), F3(method) 해소
2. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
3. Q1 해소 (pkg/queue — nesting depth 4)
4. `go build ./...` + `go test ./...`
5. `filefunc validate` pkg/ 단위 확인

## .ffignore

```
artifacts/
specs/
```

## 검증

1. `filefunc validate` — pkg/ 관련 ERROR 0
2. `go test ./...` — 전체 통과
