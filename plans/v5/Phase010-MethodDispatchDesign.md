# Phase010 Part B1 — MethodDispatchDesign

> `pkg/generate/gogin/generate_method_from_iface.go` 의 7-case switch 를 **순수 `Decide*` 함수로 수렴**하는 설계.
> 기준점 (2-depth) 재평가 반영 — Toulmin graph 불채택.

---

## 1. 현행 구조 분석

### 1.1 대상 함수

`generateMethodFromIface(b, implName, modelName, m, query, seqType, table, includes, cursorSpecs)`
— `pkg/generate/gogin/generate_method_from_iface.go:1~75`
— 호출자: `writeModelMethods()` (`write_model_methods.go:18`)

### 1.2 7-case dispatch

| # | 조건 | 이름 접두 | QueryOpts | 반환 제네릭 | 슬라이스 | seqType | 핸들러 |
|---|------|---------|-----------|------------|---------|---------|--------|
| 0 | `m.Name == "WithTx"` | `WithTx` | — | — | — | — | early return |
| 1 | `isList && isCursorReturn` | `List*` | ✓ | `*pagination.Cursor[T]` | ✓ | — | `writeCursorPaginationMethod` |
| 2 | `isList` | `List*` | ✓ | `*pagination.Page[T]` / `[]T` | ✓ | — | `writeOffsetPaginationMethod` |
| 3 | `isSliceReturn` | `*By*` | — | `[]T` | ✓ | — | `writeSliceReturnMethod` |
| 4 | `isFind \|\| seqType=="get"` | `Find*`/`Get*` | — | `*T` | — | `get` | `writeFindMethod` |
| 5 | `seqType == "post"` | `Create*` | — | `*T` | — | `post` | inline QueryRowContext |
| 6 | `seqType=="put" \|\| "delete"` | `Update*`/`Delete*` | — | `error` | — | `put`/`delete` | inline ExecContext |
| 7 | default | — | — | cardinality 기반 | — | — | fallback (`if Cardinality=="one"` → QueryRow else Exec) |

### 1.3 depth 평가 (기준점 적용)

현행 구조를 `switch` 단일 문으로 재표현:

```go
switch {
case facts.MethodName == "WithTx":                    // depth 1
    return PatternSkip
case facts.IsListPrefix && facts.IsCursorReturn:      // depth 1 (AND 은 수평)
    return PatternCursorPagination
case facts.IsListPrefix:                              // depth 1
    return PatternOffsetPagination
case facts.IsSliceReturn:                             // depth 1
    return PatternSliceReturn
case facts.IsFindPrefix || facts.SeqType == "get":    // depth 1
    return PatternFind
case facts.SeqType == "post":                         // depth 1
    return PatternCreate
case facts.SeqType == "put" || facts.SeqType == "delete":  // depth 1
    return PatternUpdateDelete
default:
    if facts.Cardinality == "one" {                   // depth 2
        return PatternFallbackOne
    }
    return PatternFallbackExec
}
```

**최대 depth = 2** (default case 내부 if). **기준점 이내** → **Toulmin 제외**.

---

## 2. 설계

### 2.1 반환 타입

```go
// pkg/generate/gogin/decide_method_pattern.go
type Pattern int

const (
    PatternSkip Pattern = iota              // WithTx
    PatternCursorPagination
    PatternOffsetPagination
    PatternSliceReturn
    PatternFind                             // QueryRowContext + scanT
    PatternCreate                           // POST 인라인
    PatternUpdateDelete                     // PUT/DELETE 인라인
    PatternFallbackOne                      // Cardinality=="one"
    PatternFallbackExec                     // Cardinality 외
)

type MethodFacts struct {
    MethodName       string
    ReturnSig        string
    SeqType          string   // "get"/"post"/"put"/"delete"/""
    Cardinality      string   // "one"/"many"/""
    HasQueryOpts     bool
    IsListPrefix     bool     // isListMethod(Name) && HasQueryOpts
    IsFindPrefix     bool     // Name startsWith "Find"/"Get"
    IsCursorReturn   bool     // contains "pagination.Cursor["
    IsSliceReturn    bool     // starts with "[]"
}
```

`MethodFacts` 채우기는 **judge 함수 밖**(호출자) 또는 **`NewMethodFacts(m, query, seqType)` 헬퍼**가 담당. judge 자체는 bool/string 필드만 참조 → 단위 테스트 용이.

### 2.2 판정 함수

```go
func DecideMethodPattern(facts MethodFacts) Pattern {
    switch {
    case facts.MethodName == "WithTx":
        return PatternSkip
    case facts.IsListPrefix && facts.IsCursorReturn:
        return PatternCursorPagination
    case facts.IsListPrefix:
        return PatternOffsetPagination
    case facts.IsSliceReturn:
        return PatternSliceReturn
    case facts.IsFindPrefix || facts.SeqType == "get":
        return PatternFind
    case facts.SeqType == "post":
        return PatternCreate
    case facts.SeqType == "put" || facts.SeqType == "delete":
        return PatternUpdateDelete
    default:
        if facts.Cardinality == "one" {
            return PatternFallbackOne
        }
        return PatternFallbackExec
    }
}
```

### 2.3 파일 배치

```
pkg/generate/gogin/
├── decide_method_pattern.go         신설 — Pattern enum + MethodFacts + DecideMethodPattern + NewMethodFacts
├── decide_method_pattern_test.go    신설 — 9 Pattern 전수 테이블 테스트
└── generate_method_from_iface.go    수정 — switch 를 Pattern 소비로 교체
```

### 2.4 호출자 소비 패턴

```go
func generateMethodFromIface(b *strings.Builder, implName, modelName string,
    m ifaceMethod, query *sqlcQuery, seqType string, table *ddlTable,
    includes []includeMapping, cursorSpecs map[string]string) {

    facts := NewMethodFacts(m, query, seqType)
    switch DecideMethodPattern(facts) {
    case PatternSkip:
        return
    case PatternCursorPagination:
        writeCursorPaginationMethod(...)
    case PatternOffsetPagination:
        writeOffsetPaginationMethod(...)
    case PatternSliceReturn:
        writeSliceReturnMethod(...)
    case PatternFind:
        writeFindMethod(...)
    case PatternCreate:
        // inline L50-53 as-is
    case PatternUpdateDelete:
        // inline L55-59 as-is
    case PatternFallbackOne:
        // inline L63-67 as-is
    case PatternFallbackExec:
        // inline L68-71 as-is
    }
}
```

호출자에 **판정 로직이 남지 않는다**. 순수 dispatcher.

---

## 3. 검증

- `go test ./pkg/generate/gogin/...` — `DecideMethodPattern` 테이블 테스트 (9 Pattern)
- `fullend gen dummys/gigbridge/specs /tmp/x` — 생성된 `internal/model/*.go` 가 기존과 실질 동일
- `cd /tmp/x/backend && go build ./...` 성공

### 3.1 동일성 보증

`classify_method_pattern_test.go` 에 **현행 로직 레퍼런스**를 테스트 내부에 복제, `DecideMethodPattern` 출력과 임의 `MethodFacts` 조합에 대해 `require.Equal` 비교.

---

## 4. 보류

- `generate_method_from_iface.go` 의 handler 들(`writeCursor*` 등) 내부 리팩토링은 Phase010 범위 밖.
- 새 Pattern 추가 시 enum + switch case 한 줄씩 추가.
- 만약 장래 이 dispatch 에 **규칙 기반 policy**(사용자 설정으로 Pattern 선택 우선순위 변경 등) 요구가 생기면 그때 Toulmin 으로 승격.
