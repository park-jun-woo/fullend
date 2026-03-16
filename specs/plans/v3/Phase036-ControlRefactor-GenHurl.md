# Phase036: 1-func-1-control 리팩토링 + Q1/Q3 — gen/hurl

## 목표

`internal/gen/hurl/`의 filefunc 위반 **26건 → 0건**. 1-func-1-control + Q1 ≤ 2 + Q3 control별 상한.

## 현황 (Phase033 이후)

- Q1 (depth > 2): 25건
- Q3 (줄 수 초과): 1건 — `buildScenarioOrder`(130줄)
- mixed control: 0건

## control 분류 영향

| 함수 | 줄 | control | Q3 상한 | 상태 |
|---|---|---|---|---|
| `buildScenarioOrder` | 130 | **iteration** | 100 | ❌ 초과 |

## 수정 계획

### Q3 분리

| 함수 | 줄 | 수정 |
|---|---|---|
| `buildScenarioOrder` | 130 | 3블록 추출: `collectAllSteps(doc)`, `classifyMidStep(s, stateOps, ...)` (selection), `filterPrereqSteps(midSteps, authFKPrefixes)` |

### Q1 해결

**extract-func (10건):**
- `buildBranchSkipSet` (depth 4) — `findBestEvent(events, transOrder)` 추출
- `buildTransitionOrder` (depth 5) — `bfsTransitions(d, order, &idx)` 추출
- `buildResourceFirstTransition` (depth 5) — `updateResourceOrder(...)` 추출
- `findTokenJSONPath` (depth 4) — `findNestedTokenField(name, prop)` 추출
- `parseDataCheckEnums` (depth 3) — `extractCheckEnums(content, re, valRe)` 추출
- `parseDDLFilesHurl` (depth 3) — `extractFKReferences(content, fkRe)` 추출
- `generateRequestBodyWithOverrides` (depth 3) — `resolveFieldLine(...)` 추출
- `writeAuthSection` (depth 3) — `findAuthOps(doc)` 추출
- `writeAuthPair` (depth 3) — `writeAsserts(buf, asserts)` 추출
- `collectAuthFKResources` (depth 3) — `extractFKPrefixes(schema)` 추출

**merge-conditions (8건):**
- `findMatchingCapture`, `findParentResource`, `generateHurlTests`, `generateResponseAssertions`, `inferCaptureField`, `pathParamUtil`, `substitutePathParams`, `buildOperationRoleMap`

**early-continue/return (4건):**
- `getResponseSchema`, `getSuccessHTTPCode`, `sortDeletesByFK`, `topoSortDelete`

**stdlib 교체 (1건):**
- `sortStringSlice` — 버블소트 → `sort.Strings()`

## 추출 예상 신규 파일: ~12개

## 검증

1. `go build ./internal/gen/hurl/`
2. `filefunc validate` — gen/hurl Q1=0, Q3=0
3. `fullend gen` → smoke.hurl 생성 결과 비교
