# Phase039: filefunc 전체 적용 — 미관리 패키지 전수 전환

## 목표

`.ffignore` 편법 제거 후 **ERROR 301 + WARNING 14 → 0건**. 전체 코드베이스 filefunc 준수.

## 현황

| 위반 유형 | 건수 |
|---|---|
| A1 | 84 |
| A3 | 84 |
| F1 | 43 |
| F2 | 37 |
| F3 | 11 |
| Q1 | 42 |
| Q3 | 14 (WARNING) |

미관리 파일: **86개** → 분해 후 예상 **~564개** (478func + 86type)

## 규모 비교

| 구분 | Phase031~038 | Phase039 |
|---|---|---|
| Before 파일 | 87 | 86 |
| After 파일 (예상) | 353 | ~564 |
| 100줄+ 분해 대상 | 5개 | 44개 |
| 최대 파일 | 1,113줄 | 752줄 |

Phase031~038과 동일 규모. 4 step으로 분할 실행.

## Step 1: Tier 3 + Tier 4 (소형 패키지, 빠른 완료)

**Tier 3 (internal 소형, 10개 패키지, 17파일)**

| 패키지 | 파일 | func | type | 분해 후 | feature |
|---|---|---|---|---|---|
| policy | 2 | 8 | 3 | ~11 | policy |
| statemachine | 2 | 6 | 2 | ~8 | statemachine |
| funcspec | 1 | 6 | 2 | ~8 | funcspec |
| genmodel | 1 | 20 | 4 | ~24 | genmodel |
| genapi | 1 | 0 | 4 | ~4 | genapi |
| projectconfig | 1 | 2 | 13 | ~15 | projectconfig |
| gen/react | 1 | 15 | 0 | ~15 | gen-react |
| cmd/fullend | 1 | 12 | 0 | ~12 | cli |
| reporter | 4 | 5 | 4 | ~9 | reporter |
| scenario | 0 | 0 | 0 | 0 | scenario |

**Tier 4 (pkg/ 유틸리티, 13개 패키지, 26파일)**

| 패키지 | 파일 | func | type | 분해 후 |
|---|---|---|---|---|
| pkg/auth | 3 | 3 | 6 | ~9 |
| pkg/authz | 1 | 3 | 3 | ~6 |
| pkg/cache | 1 | 8 | 4 | ~12 |
| pkg/config | 1 | 2 | 0 | ~2 |
| pkg/crypto | 4 | 4 | 8 | ~12 |
| pkg/file | 1 | 8 | 3 | ~11 |
| pkg/image | 2 | 2 | 4 | ~6 |
| pkg/mail | 2 | 2 | 4 | ~6 |
| pkg/pagination | 2 | 0 | 2 | ~2 |
| pkg/queue | 1 | 8 | 2 | ~10 |
| pkg/session | 1 | 8 | 4 | ~12 |
| pkg/storage | 4 | 4 | 6 | ~10 |
| pkg/text | 3 | 3 | 6 | ~9 |

Step 1 합계: 43파일 → ~213파일

## Step 2: Tier 2 (중간 패키지)

| 패키지 | 파일 | func | type | 분해 후 | 최대 파일 | feature |
|---|---|---|---|---|---|---|
| stml/generator | 5 | 48 | 5 | ~53 | 516줄 | stml-gen |
| stml/validator | 3 | 25 | 19 | ~44 | 284줄 | stml-validate |
| stml/parser | 2 | 28 | 12 | ~40 | 705줄 | stml-parse |
| ssac/parser | 2 | 26 | 8 | ~34 | 752줄 | ssac-parse |
| contract | 5 | 21 | 6 | ~27 | 290줄 | contract |

Step 2 합계: 17파일 → ~198파일

## Step 3: Tier 1-a (crosscheck)

| 패키지 | 파일 | func | type | 분해 후 | 최대 파일 | feature |
|---|---|---|---|---|---|---|
| crosscheck | 19 | 76 | 5 | ~81 | 553줄 | crosscheck |

## Step 4: Tier 1-b (ssac/generator)

| 패키지 | 파일 | func | type | 분해 후 | 최대 파일 | feature |
|---|---|---|---|---|---|---|
| ssac/generator | 10 | 62 | 10 | ~72 | 487줄 | ssac-gen |

## 실행 절차 (각 Step 공통)

1. 파일 분해 — F1(1func/file), F2(1type/file), F3(method) 해소
2. `//ff:func`/`//ff:type`/`//ff:what` 부착 — A1/A3 해소
3. `control=`/`dimension=` 부착 — A9/A15 해소
4. early-continue/추출 — Q1 해소
5. `go build ./...` + `go test ./...`
6. `filefunc validate` 패키지 단위 확인

## .ffignore

```
specs/
artifacts/
```

## 검증

1. `filefunc validate` — ERROR 0, WARNING 0
2. `go test ./...` — 전체 통과
3. `fullend validate specs/dummys/gigbridge-try02/`
4. `fullend validate specs/dummys/zenflow-try05/`
