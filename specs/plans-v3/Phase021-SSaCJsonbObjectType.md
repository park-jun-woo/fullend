# Phase021: JSONB/object 타입 매핑 수정 (SSaC 코드젠 + gluegen 모델) ✅ 완료

## 목표

PostgreSQL `JSONB` 타입과 OpenAPI `type: object` / `type: array`를 `json.RawMessage`로 올바르게 매핑하여, 생성된 핸들러·모델·struct가 JSONB 컬럼을 정상 처리하도록 수정한다.

## 근본 원인

### 데이터 흐름 (3곳 동일 문제)

```
[SSaC 코드젠 — 핸들러 생성]
DDL: payload_template JSONB
  → pgTypeToGo("JSONB") = "string"  ← 문제 ①
  → var req struct { PayloadTemplate string }

OpenAPI: payload_template: { type: object }
  → oaTypeToGo("object", "") = "string"  ← 문제 ②

[gluegen — 모델 struct + impl 생성]
DDL: payload_template JSONB
  → sqlTypeToGo("JSONB") = "string"  ← 문제 ③
  → type Action struct { PayloadTemplate string }
  → func Create(..., payloadTemplate string, ...)
```

### 서버 에러 경로

```
Hurl: {"payload_template": {}, ...}
→ ShouldBindJSON → PayloadTemplate는 string인데 {} (object)는 바인딩 불가
→ 400 "invalid request body"
```

### 문제의 본질

`pgTypeToGo()`, `oaTypeToGo()`, `sqlTypeToGo()` 세 함수 모두 JSONB/object/array를 처리하지 못해 default(`"string"`)로 빠진다.

## 영향 범위

- SSaC 핸들러 코드젠: JSONB 컬럼이 있는 request body (`var req struct`)
- gluegen 모델 코드젠: JSONB 컬럼이 있는 model struct + method signature
- 현재 zenflow의 `payload_template` (JSONB)이 해당

## 수정 포인트

### 1. `pgTypeToGo()` — JSONB/JSON 케이스 추가

파일: `internal/ssac/validator/symbol.go:1058`

```go
case "JSONB", "JSON":
    return "json.RawMessage"
```

### 2. `oaTypeToGo()` — object/array 케이스 추가

파일: `internal/ssac/validator/symbol.go:1080`

```go
case "object", "array":
    return "json.RawMessage"
```

### 3. `buildJSONBodyParams()` + `collectImports()` — encoding/json import 추가

#### 3a. `buildJSONBodyParams()` — json.RawMessage 파라미터 전달

파일: `internal/ssac/generator/go_params.go:113`

현재 `time.Time`만 import 감지용으로 전달하고 있다 (line 113-118).
`json.RawMessage`도 동일하게 전달해야 `collectImports()`에서 감지 가능:

기존:
```go
for _, rp := range rawParams {
    if rp.goType == "time.Time" {
        result = append(result, typedRequestParam{name: rp.name, goType: rp.goType})
        break
    }
}
```

수정:
```go
for _, rp := range rawParams {
    if rp.goType == "time.Time" || rp.goType == "json.RawMessage" {
        result = append(result, typedRequestParam{name: rp.name, goType: rp.goType})
    }
}
```

> 참고: 기존 `time.Time`에서 `break`를 하고 있지만, 여러 타입의 import를 동시에 감지해야 하므로 `break` 제거.

#### 3b. `collectImports()` — json.RawMessage 케이스 추가

파일: `internal/ssac/generator/go_params.go:218`

기존 goType switch:
```go
for _, tp := range reqParams {
    switch tp.goType {
    case "int64", "float64", "bool":
        seen["strconv"] = true
    case "time.Time":
        seen["time"] = true
    }
}
```

추가:
```go
    case "json.RawMessage":
        seen["encoding/json"] = true
```

`order` 배열 (line 246)에도 추가:
```go
order := []string{"database/sql", "encoding/json", "net/http", "strconv", "time"}
```

### 4a. `parseDDLFiles()` — colRe 정규식에 JSONB/JSON 추가

파일: `internal/gluegen/model_impl.go:726`

기존 `colRe`에 `JSONB|JSON`이 없어서 JSONB 컬럼이 아예 파싱되지 않았다:
```go
colRe := regexp.MustCompile(`...|TIMESTAMPTZ|TIMESTAMP|JSONB|JSON)`)
```

### 4b. `sqlTypeToGo()` — JSONB/JSON 케이스 추가

파일: `internal/gluegen/model_impl.go:905`

```go
case "JSONB", "JSON":
    return "json.RawMessage"
```

### 5. `generateTypesFile()` — encoding/json import 추가

파일: `internal/gluegen/model_impl.go:184`

기존 `needsTime` 패턴과 동일하게 `needsJSON` 추가:
```go
needsJSON := false
for _, m := range models {
    t := tables[m]
    if t == nil { continue }
    for _, col := range t.Columns {
        if col.GoType == "json.RawMessage" {
            needsJSON = true
            break
        }
    }
    if needsJSON { break }
}
```

import 출력 (기존 `needsTime` 단독 → 복수 import 블록으로 변경):
```go
var imports []string
if needsJSON {
    imports = append(imports, "\"encoding/json\"")
}
if needsTime {
    imports = append(imports, "\"time\"")
}
if len(imports) > 0 {
    b.WriteString("import (\n")
    for _, imp := range imports {
        b.WriteString("\t" + imp + "\n")
    }
    b.WriteString(")\n\n")
}
```

### 6. `renderInterfaces()` — models_gen.go encoding/json import 추가

파일: `internal/ssac/generator/go_interface.go:354`

`pgTypeToGo()` 수정 후 DDL JSONB 컬럼의 GoType이 `json.RawMessage`로 바뀌면,
`resolveInputParamType()` → `st.DDLTables[].Columns[]` 경로로 `models_gen.go` 메서드 시그니처에 `json.RawMessage`가 들어간다.
기존 `needsTimeImport()` 패턴과 동일하게 `needsJSONImport()` 추가:

```go
func needsJSONImport(interfaces []derivedInterface) bool {
    for _, iface := range interfaces {
        for _, m := range iface.Methods {
            for _, p := range m.Params {
                if p.GoType == "json.RawMessage" {
                    return true
                }
            }
        }
    }
    return false
}
```

`renderInterfaces()` import 블록 (line 360 부근)에 추가:
```go
needJSON := needsJSONImport(interfaces)
// ...
if needJSON {
    buf.WriteString("\t\"encoding/json\"\n")
}
```

### 7. `generateModelFile()` — per-model impl 파일 encoding/json import 추가

파일: `internal/gluegen/model_impl.go:240`

`generateModelFile()`이 생성하는 per-model 구현 파일(예: `action.go`)의 import 블록(line 259-268)에
`json.RawMessage` 파라미터가 있을 때 `"encoding/json"` 추가:

```go
needsJSON := false
for _, method := range methods {
    for _, p := range method.Params {
        if p.Type == "json.RawMessage" {
            needsJSON = true
            break
        }
    }
    if needsJSON { break }
}
// import 블록 내:
if needsJSON {
    b.WriteString("\t\"encoding/json\"\n")
}
```

## 변경 파일 목록

| 파일 | 변경 내용 |
|---|---|
| `internal/ssac/validator/symbol.go` | `pgTypeToGo`: JSONB/JSON → `json.RawMessage`, `oaTypeToGo`: object/array → `json.RawMessage` |
| `internal/ssac/generator/go_params.go` | `collectImports`: `json.RawMessage` 타입 시 `encoding/json` import 추가 |
| `internal/ssac/generator/go_interface.go` | `renderInterfaces`: `json.RawMessage` 파라미터 시 `encoding/json` import 추가 |
| `internal/gluegen/model_impl.go` | `sqlTypeToGo`: JSONB/JSON → `json.RawMessage`, `generateTypesFile`: `encoding/json` import 추가, `generateModelFile`: `encoding/json` import 추가 |

## 검증 방법

1. `go build ./cmd/fullend/` — fullend 빌드 통과
2. `go test ./...` — 테스트 통과
3. `fullend gen specs/zenflow artifacts/zenflow` — 재생성
4. 생성된 파일 확인:
   - `artifacts/zenflow/backend/internal/model/models_gen.go`: 메서드 시그니처에 `json.RawMessage`, import에 `"encoding/json"`
   - `artifacts/zenflow/backend/internal/model/types.go`: `PayloadTemplate json.RawMessage`, import에 `"encoding/json"`
   - `artifacts/zenflow/backend/internal/model/action.go`: import에 `"encoding/json"`
   - `artifacts/zenflow/backend/internal/service/action/create_action.go`: `PayloadTemplate json.RawMessage`, import에 `"encoding/json"`
5. zenflow 서버 빌드 → hurl 테스트에서 CreateAction 통과 (400이 아닌 200)
