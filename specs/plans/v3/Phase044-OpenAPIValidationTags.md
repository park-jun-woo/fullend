# Phase044: OpenAPI 기반 입력 검증 태그 코드젠 + validate 제약 누락 검출 ✅ 완료

## 목표

두 가지를 동시에 달성한다:

1. **validate**: OpenAPI requestBody에 검증 제약이 누락된 필드를 ERROR로 검출한다
2. **gen**: OpenAPI requestBody 스키마의 검증 제약을 읽어 SSaC 핸들러의 JSON body struct에 Gin `binding` 태그를 자동 생성한다

설계 원칙:
1. **SSOT 1회 정의** — 검증 규칙은 OpenAPI에만 존재, 코드에 수동 작성 금지
2. **기존 파이프라인 최소 확장** — `SymbolTable`에 OpenAPI request 메타데이터를 추가, `buildJSONBodyParams()`에서 소비
3. **Gin binding 태그 직결** — OpenAPI 제약 → `binding:"required,email,min=8"` 1:1 변환
4. **oapi-codegen 타입 불사용** — 현재 핸들러가 익명 struct를 쓰는 구조 유지 (types.gen.go는 클라이언트용)
5. **validate 선행** — 제약 누락이 ERROR면 gen 진행 자체가 차단되므로, 스펙 작성자가 검증 규칙을 빠뜨릴 수 없다

## 동기

### 문제 1: 코드젠이 검증 제약을 무시

현재 `buildJSONBodyParams()`는 `json` 태그만 생성한다:
```go
var req struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

OpenAPI에 `required: [email, password]`, `format: email`, `minLength: 8`이 정의되어 있어도 코드젠이 이를 무시한다. 결과:
- 빈 문자열, 잘못된 형식이 DB constraint까지 내려가야 에러 발생
- 클라이언트에 불친절한 500 에러 노출
- 보안 감사에서 입력 검증 부재로 지적

### 문제 2: OpenAPI에 검증 제약 자체가 없어도 통과

zenflow-try05의 OpenAPI:
```yaml
email:
  type: string    # minLength, maxLength, format 모두 없음
password:
  type: string    # minLength 없음
```

DDL에는 `VARCHAR(255)`, `CHECK (role IN ('admin','member'))`가 있는데 OpenAPI에는 대응하는 `maxLength`, `enum`이 없다. 현재 `fullend validate`가 이를 감지하지 못한다.

**제약 누락 검출을 validate ERROR로 올려야** 스펙 작성자가 OpenAPI를 완전히 작성하게 되고, 그래야 코드젠의 binding 태그가 의미를 갖는다.

## 설계

### 1단계: `openAPISchema` YAML struct 확장

현재 `openapi_schema.go`의 `openAPISchema`에는 `Type`, `Format`, `Properties`, `Ref`만 있다. 검증 제약 필드를 추가한다:

```go
type openAPISchema struct {
    Type       string                   `yaml:"type"`
    Format     string                   `yaml:"format"`
    Properties map[string]openAPISchema `yaml:"properties"`
    Ref        string                   `yaml:"$ref"`
    Required   []string                 `yaml:"required"`   // NEW
    MinLength  *int                     `yaml:"minLength"`  // NEW
    MaxLength  *int                     `yaml:"maxLength"`  // NEW
    Minimum    *float64                 `yaml:"minimum"`    // NEW
    Maximum    *float64                 `yaml:"maximum"`    // NEW
    Pattern    string                   `yaml:"pattern"`    // NEW
    Enum       []interface{}            `yaml:"enum"`       // NEW — YAML은 다양한 타입 허용
}
```

`openAPISchema`는 requestBody 스키마와 components 스키마 양쪽에서 재귀적으로 사용되므로, 이 확장으로 `$ref` 해석 시에도 제약 정보가 보존된다.

### 2단계: FieldConstraint + RequestSchema 구조체

신규 파일로 추가 (filefunc 1파일 1타입 원칙):

```go
// FieldConstraint는 OpenAPI schema property의 검증 제약을 담는다.
type FieldConstraint struct {
    Required  bool
    Format    string   // "email", "uuid", "date-time", "uri" 등
    MinLength *int
    MaxLength *int
    Minimum   *float64
    Maximum   *float64
    Pattern   string   // 정규식
    Enum      []string // 허용 값 목록
}

// RequestSchema는 하나의 operationId에 대한 requestBody 필드별 제약을 담는다.
type RequestSchema struct {
    Fields map[string]FieldConstraint // JSON field name → 제약
}
```

`SymbolTable` (`symbol_table.go:6-12`)에 필드 추가:
```go
type SymbolTable struct {
    Models         map[string]ModelSymbol
    Operations     map[string]OperationSymbol
    Funcs          map[string]bool
    DDLTables      map[string]DDLTable
    DTOs           map[string]bool
    RequestSchemas map[string]RequestSchema // operationId → RequestSchema (NEW)
}
```

`Clone()` (`symbol_table.go:17-33`)에 `RequestSchemas: st.RequestSchemas` 추가 (읽기 전용이므로 shallow copy).

### 3단계: loadOpenAPI()에서 RequestSchemas 수집

requestBody에서 RequestSchema를 수집하는 로직을 추가한다. 현재 `buildOperationSymbol()` (`symbol_table_build_operation_symbol.go:27-37`)에서 requestBody를 처리한다:

```go
// request body fields
if op.RequestBody == nil {
    return opSym
}
content, ok := op.RequestBody.Content["application/json"]
if !ok {
    return opSym
}
for _, f := range collectSchemaFields(content.Schema, schemas) {
    opSym.RequestFields[f] = true
}
return opSym
```

그러나 `buildOperationSymbol`은 `OperationSymbol`만 반환하고 `SymbolTable`을 직접 수정하지 않는다. RequestSchema 수집은 **`loadOpenAPI()` (`load_openapi.go:29-36`)** 루프 안에서 한다:

```go
for _, pathItem := range spec.Paths {
    for _, op := range pathItem.operations() {
        if op.OperationID == "" {
            continue
        }
        opSym := st.buildOperationSymbol(op, schemas)
        st.Operations[op.OperationID] = opSym

        // NEW: 검증 제약 수집
        if op.RequestBody != nil {
            if content, ok := op.RequestBody.Content["application/json"]; ok {
                rs := extractRequestSchema(content.Schema, schemas)
                st.RequestSchemas[op.OperationID] = rs
            }
        }
    }
}
```

**`resolveSchema` 신규 작성** (`resolve_schema.go`) — 기존 `collectSchemaFields()` (`collect_schema_fields.go`)는 `$ref`를 해석하여 필드 이름만 반환한다. `resolveSchema`는 스키마 자체를 반환하는 별도 함수로 새로 작성한다:

```go
// resolveSchema는 $ref가 있으면 해석하여 실제 스키마를 반환한다.
func resolveSchema(schema openAPISchema, allSchemas map[string]openAPISchema) openAPISchema {
    if schema.Ref == "" {
        return schema
    }
    name := schema.Ref[strings.LastIndex(schema.Ref, "/")+1:]
    if resolved, ok := allSchemas[name]; ok {
        return resolved
    }
    return schema
}
```

`extractRequestSchema` 헬퍼:
```go
func extractRequestSchema(schema openAPISchema, allSchemas map[string]openAPISchema) RequestSchema {
    resolved := resolveSchema(schema, allSchemas)
    requiredSet := map[string]bool{}
    for _, r := range resolved.Required {
        requiredSet[r] = true
    }
    rs := RequestSchema{Fields: map[string]FieldConstraint{}}
    for name, prop := range resolved.Properties {
        prop = resolveSchema(prop, allSchemas) // 중첩 $ref 해석
        fc := FieldConstraint{
            Required:  requiredSet[name],
            Format:    prop.Format,
            MinLength: prop.MinLength,
            MaxLength: prop.MaxLength,
            Minimum:   prop.Minimum,
            Maximum:   prop.Maximum,
            Pattern:   prop.Pattern,
        }
        for _, e := range prop.Enum {
            if s, ok := e.(string); ok {
                fc.Enum = append(fc.Enum, s)
            }
        }
        rs.Fields[name] = fc
    }
    return rs
}
```

### 4단계: DDLTable 확장 — VARCHAR 길이 + CHECK enum

현재 `DDLTable` (`ddl_table.go:6-12`):
```go
type DDLTable struct {
    Columns     map[string]string
    ColumnOrder []string
    ForeignKeys []ForeignKey
    Indexes     []Index
    PrimaryKey  []string
}
```

추가:
```go
type DDLTable struct {
    Columns     map[string]string
    ColumnOrder []string
    ForeignKeys []ForeignKey
    Indexes     []Index
    PrimaryKey  []string
    VarcharLen  map[string]int      // col → VARCHAR(N)의 N (NEW)
    CheckEnums  map[string][]string // col → CHECK IN 값 목록 (NEW)
}
```

### 5단계: parseDDLTables 확장

정규식은 **패키지 레벨 변수로 프리컴파일**한다:

```go
var (
    reVarcharLen = regexp.MustCompile(`(?i)VARCHAR\((\d+)\)`)
    reCheckEnum  = regexp.MustCompile(`(?i)CHECK\s*\(\s*(\w+)\s+IN\s*\(([^)]+)\)\s*\)`)
)
```

**VARCHAR 길이 파싱** — `parseDDLTables()` (`parse_ddl_tables.go:62-79`)의 컬럼 라인 처리에서:

현재 `pgTypeToGo(colType)`는 `VARCHAR(255)` → `string`으로 변환하며 길이를 버린다. 변경:

```go
colName := parts[0]
colType := strings.ToUpper(parts[1])
colType = strings.TrimSuffix(colType, ",")

goType := pgTypeToGo(colType)
if t, ok := tables[currentTable]; ok {
    t.Columns[colName] = goType
    t.ColumnOrder = append(t.ColumnOrder, colName)

    // NEW: VARCHAR 길이 추출
    if n := extractVarcharLen(colType); n > 0 {
        if t.VarcharLen == nil {
            t.VarcharLen = map[string]int{}
        }
        t.VarcharLen[colName] = n
    }
    // ... 기존 인라인 PK, UNIQUE, FK 처리
}
```

```go
func extractVarcharLen(colType string) int {
    m := reVarcharLen.FindStringSubmatch(colType)
    if len(m) == 2 {
        n, _ := strconv.Atoi(m[1])
        return n
    }
    return 0
}
```

주의: `colType`은 `parts[1]`만 보므로 `VARCHAR(255)`가 통째로 들어올 수도 있고 `VARCHAR`만 들어올 수도 있다. 전체 라인에서 파싱해야 할 수 있음 — 실제 DDL 파일의 `VARCHAR(255)` 형태를 확인하여 결정.

**CHECK enum 파싱** — 현재 `parse_ddl_tables.go:56-59`:

```go
// CHECK → skip
if strings.HasPrefix(upper, "CHECK") || line == "" {
    continue
}
```

변경:

```go
if strings.HasPrefix(upper, "CHECK") {
    if col, vals := parseCheckEnum(line); col != "" {
        if t, ok := tables[currentTable]; ok {
            if t.CheckEnums == nil {
                t.CheckEnums = map[string][]string{}
            }
            t.CheckEnums[col] = vals
            tables[currentTable] = t
        }
    }
    continue
}
```

```go
// parseCheckEnum은 CHECK (col IN ('a','b','c')) 형태를 파싱한다.
// 복합 CHECK 식은 미지원 — 단일 컬럼 IN 리스트만 파싱.
func parseCheckEnum(line string) (string, []string) {
    m := reCheckEnum.FindStringSubmatch(line)
    if len(m) < 3 {
        return "", nil
    }
    col := m[1]
    rawVals := m[2]
    var vals []string
    for _, v := range strings.Split(rawVals, ",") {
        v = strings.TrimSpace(v)
        v = strings.Trim(v, "'\"")
        if v != "" {
            vals = append(vals, v)
        }
    }
    return col, vals
}
```

또한 컬럼 라인에 인라인 CHECK가 있는 경우도 처리:
```sql
role VARCHAR(50) NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
```

컬럼 라인 처리 블록에서 인라인 CHECK 감지:
```go
// NEW: 인라인 CHECK enum
if strings.Contains(upper, "CHECK") {
    if _, vals := parseCheckEnum(line); len(vals) > 0 {
        if t.CheckEnums == nil {
            t.CheckEnums = map[string][]string{}
        }
        t.CheckEnums[colName] = vals
    }
}
```

### 6단계: crosscheck — OpenAPI 제약 누락 검출 (validate ERROR)

`internal/crosscheck/openapi_constraints.go` 신규 파일.

```go
//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what OpenAPI requestBody 검증 제약 누락·불일치를 DDL 기준으로 검출
func CheckOpenAPIConstraints(input *CrossValidateInput) []CrossError
```

`CrossValidateInput` (`internal/crosscheck/types.go:10-20`)에 이미 `*genapi.ParsedSSOTs`가 임베딩되어 있으므로 `input.SymbolTable.RequestSchemas`와 `input.SymbolTable.DDLTables`에 직접 접근 가능하다. **별도 필드 추가 불필요.**

**ERROR 규칙:**

| # | 규칙 | 조건 | 에러 메시지 예시 |
|---|---|---|---|
| C1 | **required 누락** | SSaC에서 `request.X`로 참조하는 필드가 OpenAPI requestBody의 `required` 배열에 없음 | `Register: field "email" used in SSaC but not marked required in OpenAPI` |
| C2 | **string 필드 maxLength 누락** | OpenAPI `type: string` 필드에 대응하는 DDL 컬럼이 `VARCHAR(N)`인데 OpenAPI에 `maxLength` 없음 | `Register: field "email" maps to VARCHAR(255) but OpenAPI has no maxLength` |
| C3 | **DDL CHECK ↔ enum 누락** | DDL에 `CHECK (col IN ('a','b'))`가 있는데 OpenAPI에 `enum` 없음 | `Register: field "role" has DDL CHECK constraint but no OpenAPI enum` |
| C4 | **DDL CHECK ↔ enum 값 불일치** | DDL CHECK 값과 OpenAPI enum 값이 다름 | `Register: field "role" OpenAPI enum [admin,user] ≠ DDL CHECK [admin,member]` |

**C1 범위 제한**: SSaC에서 `request.X`로 참조하는 필드 중 OpenAPI의 `required` 배열에 없으면 ERROR. 단, requestBody 자체에 정의된 property이면서 required에만 빠진 경우만 대상. OpenAPI에 property 정의 자체가 없는 필드는 기존 "SSaC ↔ OpenAPI" 규칙이 이미 검출한다.

**WARNING 규칙:**

| # | 규칙 | 조건 | 메시지 예시 |
|---|---|---|---|
| W1 | **maxLength > VARCHAR** | OpenAPI `maxLength`가 DDL `VARCHAR(N)`보다 큼 | `Register: field "email" maxLength(500) > VARCHAR(255), DB truncation risk` |
| W2 | **password 필드 minLength 미설정** | 필드명에 `password` 포함이고 `minLength` 없음 | `Register: field "password" has no minLength (security risk)` |
| W3 | **email 필드 format 미설정** | 필드명에 `email` 포함이고 `format: email` 없음 | `Register: field "email" has no format:email` |

DDL 컬럼 매핑은 기존 `lookupDDLType()`과 동일한 snake_case 변환으로 전체 DDL 테이블 순회.

### 7단계: rules.go에 Rule 등록

`internal/crosscheck/rules.go`의 `var rules` 슬라이스 (현재 15개 Rule)에 16번째로 추가:

```go
{
    Name: "OpenAPI Constraints", Source: "OpenAPI", Target: "DDL",
    Requires: func(in *CrossValidateInput) bool {
        return in.SymbolTable != nil && in.SymbolTable.RequestSchemas != nil && in.ServiceFuncs != nil
    },
    Check: func(in *CrossValidateInput) []CrossError {
        return CheckOpenAPIConstraints(in)
    },
},
```

기존 `Rule` 체계를 따르므로 `--skip OpenAPI Constraints`로 제외 가능.

### 8단계: FieldConstraint → binding 태그 변환

`internal/ssac/generator/build_binding_tag.go` **신규 파일**:

```go
//ff:func feature=ssac-gen type=util control=sequence
//ff:what FieldConstraint를 Gin binding 태그 문자열로 변환

func buildBindingTag(fc validator.FieldConstraint) string {
    var parts []string
    if fc.Required {
        parts = append(parts, "required")
    }
    switch fc.Format {
    case "email":
        parts = append(parts, "email")
    case "uuid":
        parts = append(parts, "uuid")
    case "uri":
        parts = append(parts, "uri")
    }
    if fc.MinLength != nil {
        parts = append(parts, fmt.Sprintf("min=%d", *fc.MinLength))
    }
    if fc.MaxLength != nil {
        parts = append(parts, fmt.Sprintf("max=%d", *fc.MaxLength))
    }
    if fc.Minimum != nil {
        parts = append(parts, fmt.Sprintf("gte=%g", *fc.Minimum))
    }
    if fc.Maximum != nil {
        parts = append(parts, fmt.Sprintf("lte=%g", *fc.Maximum))
    }
    if len(fc.Enum) > 0 {
        parts = append(parts, "oneof="+strings.Join(fc.Enum, " "))
    }
    if len(parts) == 0 {
        return ""
    }
    return `binding:"` + strings.Join(parts, ",") + `"`
}
```

`pattern`은 Gin 기본 validator에 없으므로 Phase044에서는 미지원. validate에서 `pattern` 존재 시 INFO 로그.

### 9단계: buildJSONBodyParams 수정

**대상 파일**: `internal/ssac/generator/build_json_body_params.go` (filefunc 분리 후 현행 파일)

현재 시그니처 (`build_json_body_params.go:12`):
```go
func buildJSONBodyParams(rawParams []rawParam) []typedRequestParam
```

변경:
```go
func buildJSONBodyParams(rawParams []rawParam, rs *validator.RequestSchema) []typedRequestParam
```

struct 필드 생성 부분 (`build_json_body_params.go:17`) 변경:

```go
// 현재
buf.WriteString(fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", strcase.ToGoPascal(rp.name), rp.goType, rp.name))

// 변경 후
tag := fmt.Sprintf("json:\"%s\"", rp.name)
if rs != nil {
    if fc, ok := rs.Fields[rp.name]; ok {
        if bt := buildBindingTag(fc); bt != "" {
            tag += " " + bt
        }
    }
}
buf.WriteString(fmt.Sprintf("\t\t%s %s `%s`\n", strcase.ToGoPascal(rp.name), rp.goType, tag))
```

ShouldBindJSON 에러 응답 (`build_json_body_params.go:21`) 변경:

```go
// 현재
buf.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"invalid request body\"})\n")

// 변경 후
buf.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})\n")
```

### 10단계: collectRequestParams에 operationID 전달

**대상 파일**: `internal/ssac/generator/collect_request_params.go` (filefunc 분리 후 현행 파일)

현재 시그니처 (`collect_request_params.go:11`):
```go
func collectRequestParams(seqs []parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool) []typedRequestParam
```

변경:
```go
func collectRequestParams(seqs []parser.Sequence, st *validator.SymbolTable, pathParamSet map[string]bool, operationID string) []typedRequestParam
```

내부의 `buildJSONBodyParams` 호출부 (`collect_request_params.go:15`) 변경:

```go
// 현재
return buildJSONBodyParams(rawParams)

// 변경 후
var rs *validator.RequestSchema
if st != nil && st.RequestSchemas != nil {
    if schema, ok := st.RequestSchemas[operationID]; ok {
        rs = &schema
    }
}
return buildJSONBodyParams(rawParams, rs)
```

### 11단계: 호출부 2곳 수정

`collectRequestParams`의 호출부가 **2곳** 있다. 양쪽 모두 `sf.Name` (= operationID)을 전달하도록 변경:

**호출부 1**: `go_handler.go:22` (filefunc 미적용 — `generateHTTPFunc` 내부)
```go
// 현재
requestParams := collectRequestParams(sf.Sequences, st, pathParamSet)

// 변경 후
requestParams := collectRequestParams(sf.Sequences, st, pathParamSet, sf.Name)
```

**호출부 2**: `analyze_http_func.go:17` (filefunc 적용)
```go
// 현재
requestParams := collectRequestParams(sf.Sequences, st, pathParamSet)

// 변경 후
requestParams := collectRequestParams(sf.Sequences, st, pathParamSet, sf.Name)
```

두 호출부 모두 `analyzeHTTPFunc(sf, st, g)` / `generateHTTPFunc(sf, st)` 컨텍스트에서 `sf parser.ServiceFunc`를 받으므로 `sf.Name` 접근에 문제 없다.

### 12단계: 산출물 변경 예시

**변경 전** (`register.go`):
```go
var req struct {
    Password string `json:"password"`
    OrgName  string `json:"org_name"`
    Email    string `json:"email"`
}
```

**변경 후** (OpenAPI에 제약 추가 후):
```go
var req struct {
    Password string `json:"password" binding:"required,min=8,max=128"`
    OrgName  string `json:"org_name" binding:"required,min=1,max=255"`
    Email    string `json:"email" binding:"required,email,max=255"`
}
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| **ssac/validator — symbol feature** | |
| `openapi_schema.go` | `openAPISchema`에 `Required`·`MinLength`·`MaxLength`·`Minimum`·`Maximum`·`Pattern`·`Enum` YAML 태그 추가 |
| `field_constraint.go` | **신규** — `FieldConstraint` 구조체 |
| `request_schema.go` | **신규** — `RequestSchema` 구조체 |
| `symbol_table.go` | `SymbolTable`에 `RequestSchemas` 필드 추가, `Clone()`에 shallow copy 추가 |
| `load_openapi.go` | `loadOpenAPI()` 루프에서 `extractRequestSchema()` 호출, `st.RequestSchemas` 저장 |
| `resolve_schema.go` | **신규** — `resolveSchema()`: `$ref` 해석하여 스키마 반환 |
| `extract_request_schema.go` | **신규** — `extractRequestSchema()`: requestBody 스키마에서 필드별 제약 수집 |
| `ddl_table.go` | `DDLTable`에 `VarcharLen`·`CheckEnums` 필드 추가 |
| `parse_ddl_tables.go` | CHECK enum 파싱 (기존 skip 제거), 인라인 CHECK 감지, VARCHAR 길이 추출 |
| `extract_varchar_len.go` | **신규** — `extractVarcharLen()` + 패키지 레벨 `reVarcharLen` |
| `parse_check_enum.go` | **신규** — `parseCheckEnum()` + 패키지 레벨 `reCheckEnum` |
| **crosscheck feature** | |
| `check_openapi_constraints.go` | **신규** — `CheckOpenAPIConstraints()`: C1~C4 ERROR, W1~W3 WARNING |
| `rules.go` | `rules` 슬라이스에 `"OpenAPI Constraints"` Rule 추가 (16번째) |
| **ssac/generator — ssac-gen feature** | |
| `build_binding_tag.go` | **신규** — `buildBindingTag()`: FieldConstraint → Gin binding 태그 변환 |
| `build_json_body_params.go` | 시그니처에 `*validator.RequestSchema` 추가, binding 태그 생성, 에러 메시지를 `err.Error()`로 변경 |
| `collect_request_params.go` | 시그니처에 `operationID string` 추가, `buildJSONBodyParams` 호출 시 `RequestSchema` 전달 |
| `go_handler.go` | `collectRequestParams()` 호출부에 `sf.Name` 전달 |
| `analyze_http_func.go` | `collectRequestParams()` 호출부에 `sf.Name` 전달 |
| **문서** | |
| `artifacts/manual-for-ai.md` | OpenAPI 검증 제약 필수 규칙(C1~C4) + binding 태그 매핑 규칙 문서화 |

## filefunc 어노테이션

신규 파일에 추가할 어노테이션 (filefunc 1파일 1함수/1타입 원칙 준수):

| 파일 | 어노테이션 |
|---|---|
| `field_constraint.go` | `//ff:type feature=symbol type=model` / `//ff:what OpenAPI schema property의 검증 제약` |
| `request_schema.go` | `//ff:type feature=symbol type=model` / `//ff:what operationId별 requestBody 필드 제약` |
| `resolve_schema.go` | `//ff:func feature=symbol type=util control=selection` / `//ff:what $ref가 있으면 해석하여 실제 스키마를 반환` |
| `extract_request_schema.go` | `//ff:func feature=symbol type=util control=iteration dimension=1` / `//ff:what requestBody 스키마에서 필드별 검증 제약을 수집` |
| `extract_varchar_len.go` | `//ff:func feature=symbol type=util control=sequence` / `//ff:what VARCHAR(N) 타입에서 길이 N을 추출` |
| `parse_check_enum.go` | `//ff:func feature=symbol type=parser control=sequence` / `//ff:what CHECK (col IN (...)) 절에서 컬럼명과 허용 값을 파싱` |
| `check_openapi_constraints.go` | `//ff:func feature=crosscheck type=rule control=iteration dimension=1` / `//ff:what OpenAPI requestBody 검증 제약 누락·불일치를 DDL 기준으로 검출` |
| `build_binding_tag.go` | `//ff:func feature=ssac-gen type=util control=sequence` / `//ff:what FieldConstraint를 Gin binding 태그 문자열로 변환` |

## 의존성

- Phase043(filefunc ssac/generator 적용) 완료 후. `collect_request_params.go`, `build_json_body_params.go`, `rawParam` 타입에 의존.
- 외부 패키지 추가 없음 — 자체 YAML struct (`openAPISchema`), Gin `binding` (기존).

## 검증

1. `go test ./internal/ssac/validator/...` — `DDLTable.VarcharLen`, `CheckEnums` 파싱 테스트, `RequestSchema` 수집 테스트
2. `go test ./internal/crosscheck/...` — C1~C4 ERROR, W1~W3 WARNING 규칙 테스트
3. `go test ./internal/ssac/generator/...` — `buildBindingTag()` 단위 테스트: required, email, min/max, enum, 복합 조합
4. zenflow-try05 OpenAPI에 제약 **없이** `fullend validate` → C1~C4 ERROR 발생 확인
5. zenflow-try05 OpenAPI에 제약 추가 후 `fullend validate` → 통과
6. `fullend gen` → 생성된 핸들러에 binding 태그 존재 확인
7. 생성된 프로젝트 `go build` 통과
8. Hurl 테스트: 빈 email로 register 시 400 반환
9. `go vet ./...` 통과

## 리스크

- **Gin binding 태그 한계** — `pattern` (정규식)은 Gin 기본 validator에 없어 커스텀 등록 필요. Phase044에서는 pattern을 제외하고 향후 확장.
- **기존 산출물 변경** — binding 태그 추가로 기존보다 strict해짐. OpenAPI에 제약이 없으면 validate ERROR로 차단되므로 반드시 추가해야 함.
- **DDL CHECK 파싱 복잡도** — `CHECK (col IN ('a','b'))` 외 복합 CHECK 식은 미지원. 단일 컬럼 IN 리스트만 파싱. 인라인 CHECK (컬럼 라인에 포함된 경우)도 처리.
- **VARCHAR 길이 파싱** — `parts[1]`이 `VARCHAR(255)` 통째로 올 수 있고 `VARCHAR`만 올 수도 있음. 정규식으로 전체 라인에서 추출하는 것이 안전.
- **에러 메시지 노출** — `err.Error()`가 Go struct 필드명을 노출. 프로덕션에서는 필드명 매핑 필요할 수 있으나, 코드젠 범위에서는 기본 메시지로 충분.
- **`date-time` format** — OpenAPI `format: date-time`은 Go `time.Time` 타입과 대응하지만 binding 태그로는 검증 불가. JSON unmarshal 시점에서 자연 검증되므로 별도 처리 불필요.
- **`go_handler.go` filefunc 미적용** — `go_handler.go`는 아직 `//ff:func` 어노테이션이 없는 구버전 파일. `analyze_http_func.go`와 중복 로직이 있으나, Phase044 범위에서는 호출부만 수정하고 중복 제거는 별도 Phase에서 처리.
