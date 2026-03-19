# Phase034: dimension 부착 + Q1/Q3 리팩토링 — ssac/validator ✅ 완료

## 목표

`internal/ssac/validator/`의 filefunc 위반 **0건** 달성. `dimension=` 어노테이션 부착 + early-continue 적용 + Q3 초과 분리.

## 현황 (Phase033 이후)

- A15 (`dimension=` 누락): 29건 — 모든 `control=iteration` 파일
- Q1 (depth 초과): 23건 — 대부분 dimension + early-continue로 해소 가능
- Q3 (줄 수 초과): 3건 — `loadPackageGoInterfaces`(157줄), `validateRequiredFields`(119줄), `parseDDLTables`(108줄)
- 1-func-1-control 위반: 1건 — `validateRequiredFields` (for→switch)

## 설계

### 1단계: `dimension=` 어노테이션 부착

29개 `control=iteration` 파일에 `dimension=` 추가. dimension = 함수 내 for-range 체인의 최대 중첩 수.

dimension ≥ 2 파일 (다차원 데이터 순회):

| 함수 | dimension | 순회 대상 | Q1 상한 |
|---|---|---|---|
| `loadPackageGoInterfaces` | 5 | entries→decls→specs→fields→names | 6 |
| `loadGoInterfaces` | 3 | entries→decls→specs | 4 |
| `loadOpenAPI` | 2 | paths→operations | 3 |
| `validateModel` | 2 | sequences→methods | 3 |
| `validateFKReferenceGuard` | 2 | sequences→steps | 3 |
| `hasStructField` | 2 | structs→fields | 3 |
| `validateGoReservedWords` | 2 | sequences→inputs | 3 |
| `loadSqlcQueries` | 2 | entries→queries | 3 |
| `collectSchemaFields` | 2 | properties→nested | 3 |
| `findColumnTable` | 2 | tables→columns | 3 |
| `loadPackageInterfaces` | 2 | entries→models | 3 |
| `validateCallInputTypes` | 2 | sequences→inputs | 3 |
| `validateQueryUsage` | 2 | sequences→queries | 3 |

나머지 16개 파일: `dimension=1` (flat list, Q1 상한 2)

### 2단계: early-continue 적용

`if ok { ... }` 패턴을 `if !ok { continue }` 패턴으로 변환. 파일 분리 없이 depth를 줄인다.

Q1 위반 해소 대상:

| 함수 | depth (현재) | depth (early-continue 후) | dimension | Q1 상한 | 결과 |
|---|---|---|---|---|---|
| `loadPackageGoInterfaces` | 8 | 6 | 5 | 6 | ✅ 해소 |
| `loadGoInterfaces` | 6 | 3 | 3 | 4 | ✅ 해소 |
| `loadOpenAPI` | 5 | 3 | 2 | 3 | ✅ 해소 |
| `validateModel` | 5 | 3 | 2 | 3 | ✅ 해소 |
| `validateFKReferenceGuard` | 5 | 3 | 2 | 3 | ✅ 해소 |
| `hasStructField` | 4 | 2 | 2 | 3 | ✅ 해소 |
| `loadSqlcQueries` | 4 | 2 | 2 | 3 | ✅ 해소 |
| `validateGoReservedWords` | 4 | 2 | 2 | 3 | ✅ 해소 |
| `validatePaginationType` | 4 | 3 | 1 | 2 | ❌ 3단계 |
| `validateRequest` | 4 | 2 | 1 | 2 | ✅ 해소 |
| `validateStaleResponse` | 4 | 2 | 1 | 2 | ✅ 해소 |
| `parseDDLTables` | 4 | 2 | 1 | 2 | ✅ 해소 |
| `toSnakeCase` | 4 | 2 | 1 | 2 | ✅ 해소 |
| `parseInlineFk` | 3 | 2 | 1 | 2 | ✅ 해소 |
| `resolveCallInputType` | 3 | 2 | 1 | 2 | ✅ 해소 |
| `collectSchemaFields` | 3 | 2 | 2 | 3 | ✅ 해소 |
| `findColumnTable` | 3 | 2 | 2 | 3 | ✅ 해소 |
| `loadPackageInterfaces` | 3 | 2 | 2 | 3 | ✅ 해소 |
| `validateCallInputTypes` | 3 | 2 | 2 | 3 | ✅ 해소 |
| `validateQueryUsage` | 3 | 2 | 2 | 3 | ✅ 해소 |
| `validateVariableFlow` | 3 | 2 | 1 | 2 | ✅ 해소 |
| `validateRequiredFields` | 3 | — | — | — | 3단계에서 처리 |

**early-continue로 Q1 위반 22건 → 2건.** (`validatePaginationType`, `validateRequiredFields`는 3단계에서 처리)

> `validateSubscribeRules`(depth 6)는 `control=sequence`. Q1 상한 2. 별도 처리 필요 — `validateMessageFields()` 추출.

### 3단계: Q3 초과 + 1-func-1-control 분리

2단계 후 남은 위반 (Q1 잔여 + Q3 + 1-func-1-control):

| 함수 | 위반 | 수정 |
|---|---|---|
| `loadPackageGoInterfaces` (157줄) | Q3 | 3단계 추출: `collectRequestStructs()`, `parsePackageInterfaces()`, `parseStandaloneFuncs()`. 각 ~50줄 |
| `parseDDLTables` (108줄) | Q3 | `parseColumnLine()` 추출. 본체 ~60줄 |
| `validateRequiredFields` (119줄) | Q3 + 1-func-1-control | `validateSeqRequiredFields(seq, ctx)` (selection) 추출. 원본은 for 래퍼 (iteration). 2파일, Q3 모두 OK |
| `validatePaginationType` (depth 3) | Q1 (dimension=1) | switch 본체를 `checkPaginationStyle()` 추출 |
| `validateSubscribeRules` (depth 6) | Q1 (sequence) | `validateMessageFields()` 추출 |

## 변경 파일

- `dimension=` 추가: 29개 iteration 파일
- early-continue 적용: ~22개 파일
- 함수 추출 (신규 파일): ~7개

## 검증

1. `go test ./internal/ssac/validator/...`
2. `filefunc validate` — ssac/validator Q1=0, Q3=0, A9~A16=0
3. `fullend validate specs/dummys/gigbridge-try02/`
4. mutest SSaC 49건 재실행
