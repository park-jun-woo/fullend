# Phase035: 1-func-1-control 리팩토링 + Q1/Q3 — gen/gogin

## 목표

`internal/gen/gogin/`의 filefunc 위반 **41건 → 0건**. mixed 2건 분리 + 1-func-1-control + Q1 ≤ 2 + Q3 control별 상한.

## 현황 (Phase033 이후)

- Q1 (depth > 2): 33건
- Q3 (줄 수 초과): 8건
- mixed control: 2건 (`generateMethodFromIface`, `generateStateMachineSource`)

## control 분류 영향

| 함수 | 줄 | control | Q3 상한 | 상태 |
|---|---|---|---|---|
| `generateMethodFromIface` | 285 | **selection** (switch 7 case가 본체) | 300 | ✅ OK (mixed 분리 후) |
| `generateQueryOpts` | 216 | **sequence** (백틱 할당+writeFile) | 100 | ❌ 초과. 템플릿이 212줄이지만 Q3는 총줄 기준. const 추출 필요 |
| `generateMainWithDomains` | 192 | **iteration** (for 루프가 주 제어) | 100 | ❌ 초과 |
| `generateServerStruct` | 137 | **iteration** | 100 | ❌ 초과 |
| `generateCentralServer` | 137 | **iteration** | 100 | ❌ 초과 |
| `generateMain` | 126 | **sequence** | 100 | ❌ 초과 |
| `generateModelFile` | 104 | **iteration** | 100 | ❌ 초과 |
| `transformSource` | 103 | **iteration** | 100 | ❌ 초과 |

## 수정 계획

### mixed control 분리 (2건)

| 함수 | 패턴 | 추출 |
|---|---|---|
| `generateMethodFromIface` | for(param reorder) + switch(7 case) | param reorder를 `reorderCallArgs(m, query) []string` (iteration)으로 추출. 본체는 switch만 남아 `control=selection`, 285→~250줄 (300 이내) |
| `generateStateMachineSource` | for(transitions) + switch-like(if chain) | transition 루프를 별도 함수로 추출 |

### Q3 초과 함수 분리

| 함수 | 줄 | control | 수정 |
|---|---|---|---|
| `generateQueryOpts` | 216 | sequence | 템플릿을 `const queryOptsTmpl` 파일로 추출. 함수는 `writeFile(tmpl)` ~10줄 |
| `generateMainWithDomains` | 192 | iteration | domain init builder, import builder, queue builder 3개 추출 |
| `generateServerStruct` | 137 | iteration | `collectPathParams(pathItem, op)` 추출 |
| `generateCentralServer` | 137 | iteration | `writeRoutes(b, doc, opDomains)` 추출 |
| `generateMain` | 126 | sequence | queue 블록, authz 블록 추출 |
| `generateModelFile` | 104 | iteration | `detectModelImports(methods)` 추출 |
| `transformSource` | 103 | iteration | `injectTypeAssertions()`, `fixImports()` 추출 |

### Q1 해결

**depth 4+ (extract-func):**
- `generateServerStruct` (depth 6) — collectPathParams로 동시 해결
- `collectModelIncludes` (depth 5) — dedup 루프 → set 교체
- 14개 depth 4 함수 — 각각 내부 블록 추출

**depth 3 (merge-conditions/early-continue):**
- `collectFuncs`, `collectModels` 등 8개 — 조건 병합으로 depth 2로
- `hasSeqType` 공유 헬퍼 추출 — `hasAuthSequence`, `hasPublishSequence` 통합

## 추출 예상 신규 파일: ~25개

## 검증

1. `go test ./internal/gen/gogin/...`
2. `filefunc validate` — gen/gogin Q1=0, Q3=0, A9/A10/A11=0
3. `fullend gen specs/dummys/zenflow-try05/` — 코드젠 결과 비교
4. `fullend gen specs/dummys/gigbridge-try02/` — 코드젠 결과 비교
