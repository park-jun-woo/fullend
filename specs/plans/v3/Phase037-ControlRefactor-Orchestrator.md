# Phase037: 1-func-1-control 리팩토링 + Q1/Q3 — orchestrator

## 목표

`internal/orchestrator/`의 filefunc 위반 **23건 → 0건**. 1-func-1-control + Q1 ≤ 2 + Q3 control별 상한.

## 현황 (Phase033 이후)

- Q1 (depth > 2): 19건
- Q3 (줄 수 초과): 4건
- mixed control: 0건

## control 분류 영향

| 함수 | 줄 | control | Q3 상한 | 상태 |
|---|---|---|---|---|
| `GenWith` | 151 | **sequence** (순차 호출 체인) | 100 | ❌ 초과 |
| `ValidateWith` | 102 | **selection** (switch kind 디스패처) | 300 | ✅ OK |
| `statusCmd (Status)` | 128 | **selection** (switch kind 디스패처) | 300 | ✅ OK |
| `TestParseIdempotency` | 117 | **sequence** | 100 | ❌ 초과 (F5 테스트 파일) |

`ValidateWith`와 `Status`는 `control=selection`으로 Q3 300줄 이내 — **리팩토링 불필요**.

## 수정 계획

### Q3 분리

| 함수 | 줄 | control | 수정 |
|---|---|---|---|
| `GenWith` | 151 | sequence | `genAllSSOTs(...)` (순차 블록 추출), `restorePreserved(snap)` 추출 |
| `TestParseIdempotency` | 117 | sequence (테스트) | `assertIdempotent(t, name, parseFunc)` 테스트 헬퍼 추출 |

### Q1 해결

**고위험 (depth 5+, extract-func):**
- `checkDDLNullableColumns` (depth 7) — `checkColumnLine(...)`, `checkSentinelRecord(...)` 추출
- `traceFuncSpecs` (depth 6) — `findFuncSpecFile(...)`, `locateFuncSpecFile(...)` 추출
- `findFullendPkgRoot` (depth 5) — `isFullendGoMod(data)` 추출
- `statusCmd` (depth 5) — 각 kind별 summary 헬퍼 추출

**중위험 (depth 3-4, extract-func):**
- `traceArtifacts` (depth 4) — `traceModelMethods(sf, funcs)` 추출
- `tracePolicy` (depth 4) — `matchPolicyRule(rule, ...)` 추출
- `loadDTOTypes` (depth 4) — `parseDTOsFromFile(path)` 추출
- `checkPathParamConflicts` (depth 4) — `extractParamSegments(path)` 추출
- `genWith` (depth 3) — 위 Q3 분리로 동시 해결

**저위험 (depth 3, early-continue/merge-conditions):**
- `detectSSOTs`, `determineModulePath`, `findDDLTable`, `genFunc`, `scanFuncImports`, `validateSSaC`, `stmlMatchAttr`, `toSnakeCase`, `runContractValidate`, `checkSqlcQueryDuplicates`, `traceOpenAPI`, `traceHurlScenarios` — 12건

**공유 헬퍼:**
- `findEndpointPath(doc, opID)` — `traceOpenAPI`와 `traceHurlScenarios`가 공유

## 추출 예상 신규 파일: ~20개

## 검증

1. `go test ./internal/orchestrator/...`
2. `filefunc validate` — orchestrator Q1=0, Q3=0, A9/A10/A11=0
3. `fullend validate specs/dummys/zenflow-try05/`
4. `fullend validate specs/dummys/gigbridge-try02/`
5. mutest 교차검증 25건 재실행
