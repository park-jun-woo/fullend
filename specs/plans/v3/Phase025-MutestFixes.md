# Phase025: Mutest 미검출 2건 수정 — 상태 일관성 + DDL→OpenAPI 커버리지 ✅ 완료

## 목표

Mutation Test에서 미검출된 2건을 수정한다.

1. **MUT-04**: States 상태명 내부 일관성 검증 추가 (MEDIUM)
2. **MUT-11**: DDL→OpenAPI 역방향 property 커버리지 검증 추가 (LOW)

## 의존성

없음 — Phase024 이후 독립 실행 가능.

## 발견 경위

Mutest Phase 1 재실행 결과, 23건 중 2건 미검출 확인.

---

## 이슈 A: States 상태명 내부 일관성 미검증 (MUT-04)

**현상**: `[*] --> draft`와 `Draft --> open: PublishGig`에서 `draft`와 `Draft`가 서로 다른 상태로 파싱되지만, crosscheck가 잡지 못함.

**근원 분석**:

1. 파서 (`statemachine/parser.go:42-66`)는 케이스를 그대로 보존:
   - `[*] --> draft` → `InitialState = "draft"`, `stateSet["draft"] = true`
   - `Draft --> open: PublishGig` → `Transition{From: "Draft", To: "open"}`, `stateSet["Draft"] = true`
   - `d.States = ["Draft", "draft", "open"]` — 3개 상태로 인식

2. crosscheck (`states.go:46-90`)는 **전이 이벤트명**(PublishGig)만 SSaC/OpenAPI와 대조. **상태명 자체의 일관성**은 미검증.

3. States 파서 수준에서 잡는 것이 적절 — 다이어그램 내부의 자기모순이므로 crosscheck보다 앞 단계.

**수정 방안**: `statemachine/parser.go` Parse 함수에 상태명 case-insensitive 중복 검증 추가.

```go
// parser.go — Parse() 함수, stateSet 구축 후 (line 63 이후)
// 상태명 대소문자 일관성 검증: case-insensitive로 같은 이름이 2개 이상이면 에러
lowerMap := make(map[string]string) // lowercase → first seen form
for s := range stateSet {
    low := strings.ToLower(s)
    if prev, exists := lowerMap[low]; exists && prev != s {
        return nil, fmt.Errorf("state name conflict in %s: %q and %q differ only in case", id, prev, s)
    }
    lowerMap[low] = s
}
```

이 검증은 Parse 단계에서 에러를 반환하므로, `✓ States` 줄에서 바로 잡힌다.

**영향**: 의도적으로 대소문자가 다른 상태명을 사용하는 경우는 Mermaid 규약상 존재하지 않으므로 false positive 위험 없음.

---

## 이슈 B: DDL→OpenAPI 역방향 property 커버리지 미검증 (MUT-11)

**현상**: DDL `gigs.budget` 칼럼이 존재하지만 OpenAPI Gig schema에서 `budget` property를 삭제해도 WARNING 없음.

**근원 분석**:

1. `checkGhostProperties` (`openapi_ddl.go:209-251`)은 **OpenAPI→DDL 단방향**:
   - OpenAPI에 있는 property가 DDL에 없으면 ERROR (유령 property)
   - DDL에 있는 칼럼이 OpenAPI에 없는 경우는 미검증

2. `CheckDDLCoverage` (`ddl_coverage.go:15-52`)는 DDL→SSaC만 검증.

3. DDL→OpenAPI 역방향은 경고 수준이 적절 — `@sensitive` 칼럼이나 내부 전용 칼럼(예: `password_hash`)은 의도적으로 OpenAPI에 노출하지 않을 수 있음.

**수정 방안**:

### 1. `CheckOpenAPIDDL` 시그니처 확장

`sensitiveCols`를 직접 전달한다. `sensitivePatterns` 서브스트링 매칭 대체보다 `@sensitive` 어노테이션(설계자의 명시적 의도)이 정확하다.

```go
// openapi_ddl.go — 시그니처 변경
func CheckOpenAPIDDL(doc *openapi3.T, st *ssacvalidator.SymbolTable,
    funcs []ssacparser.ServiceFunc, sensitiveCols map[string]map[string]bool) []CrossError {
    // ...
    errs = append(errs, checkGhostProperties(doc, st)...)
    errs = append(errs, checkMissingProperties(doc, st, sensitiveCols)...)  // 추가
    return errs
}
```

### 2. `crosscheck.go:40` 호출부 수정

```go
errs = append(errs, CheckOpenAPIDDL(input.OpenAPIDoc, input.SymbolTable,
    input.ServiceFuncs, input.SensitiveCols)...)
```

### 3. `checkMissingProperties` 신규 함수

```go
// openapi_ddl.go — 신규 함수
func checkMissingProperties(doc *openapi3.T, st *ssacvalidator.SymbolTable,
    sensitiveCols map[string]map[string]bool) []CrossError {
    var errs []CrossError
    if doc.Components == nil || doc.Components.Schemas == nil {
        return errs
    }

    // FK 칼럼 스킵: x-include로 조인되는 칼럼은 OpenAPI에 없을 수 있음 (정상)
    xIncludeFields := collectXIncludeLocalFields(doc)

    for schemaName, schemaRef := range doc.Components.Schemas {
        if schemaRef == nil || schemaRef.Value == nil {
            continue
        }
        schema := schemaRef.Value
        tableName := modelToTable(schemaName)
        table, ok := st.DDLTables[tableName]
        if !ok {
            continue // @dto 등 DDL 매핑 없는 스키마 — 스킵
        }

        for colName, colType := range table.Columns {
            // @sensitive 칼럼은 의도적 비노출 — 스킵
            if sensitiveCols != nil {
                if cols, ok := sensitiveCols[tableName]; ok && cols[colName] {
                    continue
                }
            }
            // sensitivePatterns 매칭 칼럼도 스킵 (미어노테이션이지만 민감 추정)
            if matchesSensitivePattern(colName) {
                continue
            }
            // x-include FK 칼럼 스킵
            if xIncludeFields[colName] {
                continue
            }
            if _, exists := schema.Properties[colName]; !exists {
                errs = append(errs, CrossError{
                    Rule:       "DDL ↔ OpenAPI",
                    Context:    fmt.Sprintf("table %s.%s", tableName, colName),
                    Message:    fmt.Sprintf("DDL column %q (%s) — OpenAPI %s schema에 해당 property 없음", colName, colType, schemaName),
                    Level:      "WARNING",
                    Suggestion: fmt.Sprintf("OpenAPI %s schema에 %s property를 추가하거나, DDL에서 제거하세요", schemaName, colName),
                })
            }
        }
    }
    return errs
}

// matchesSensitivePattern은 칼럼명이 sensitivePatterns 중 하나와 매칭되는지 확인한다.
func matchesSensitivePattern(colName string) bool {
    lower := strings.ToLower(colName)
    for _, p := range sensitivePatterns {
        if strings.Contains(lower, p) {
            return true
        }
    }
    return false
}
```

### 설계 결정 근거

| 스킵 대상 | 방법 | 이유 |
|---|---|---|
| `@sensitive` 칼럼 | `sensitiveCols` 맵 조회 | 설계자의 명시적 의도. 가장 정확 |
| sensitivePatterns 매칭 칼럼 | `matchesSensitivePattern` | `@sensitive` 미기재된 `password_hash` 등도 false positive 방지 |
| FK 조인 칼럼 | `xIncludeFields` 재활용 | `x-include: [client_id:users.id]`로 조인되는 FK는 OpenAPI 비노출 정상 |
| `id`, `created_at` 등 | 스킵하지 않음 | OpenAPI에 이미 존재하는 것이 정상. 빠져있으면 실제 누락이므로 WARNING 유효 |

**이전 계획과의 차이**: `id`, `created_at`, `updated_at` 하드코딩 스킵 제거. 이 칼럼들은 OpenAPI에 노출되는 것이 정상이므로 누락 시 WARNING이 타당. 하드코딩은 DDL에 `updated_at`이 없는 프로젝트에서 혼란 유발.

---

## 변경 파일

### 이슈 A (상태명 일관성)
- `internal/statemachine/parser.go` (:63 이후 — case-insensitive 중복 검증 추가)
- `internal/statemachine/parser_test.go` (신규 테스트: 대소문자 충돌 시 에러 반환 확인)

### 이슈 B (DDL→OpenAPI 커버리지)
- `internal/crosscheck/openapi_ddl.go` (`CheckOpenAPIDDL` 시그니처에 `sensitiveCols` 추가, `checkMissingProperties` 신규, `matchesSensitivePattern` 신규)
- `internal/crosscheck/crosscheck.go` (:40 `CheckOpenAPIDDL` 호출부에 `input.SensitiveCols` 전달)
- `internal/crosscheck/openapi_ddl_test.go` (신규 테스트: DDL 칼럼 OpenAPI 미노출 WARNING 확인, @sensitive 스킵 확인)

## 검증 방법

1. MUT-04 재현: `Draft --> open` ← States 파싱 단계에서 에러 검출 확인
2. MUT-11 재현: Gig schema에서 `budget` 제거 ← WARNING 검출 확인
3. `fullend validate specs/gigbridge` 정상 통과 (변경 없는 원본) — `@sensitive` 칼럼 false positive 없음 확인
4. `go test ./internal/statemachine/...` 기존 + 신규 통과
5. `go test ./internal/crosscheck/...` 기존 + 신규 통과

## whyso 이력

| 파일 | 이력 |
|---|---|
| `internal/statemachine/parser.go` | 없음 (최초 생성 후 미수정) |
| `internal/crosscheck/openapi_ddl.go` | Phase014에서 `checkGhostProperties` 추가 |
| `internal/crosscheck/crosscheck.go` | Phase014에서 `Run` 함수 구조 확립 |
