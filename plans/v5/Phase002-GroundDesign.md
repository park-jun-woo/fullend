# Phase002 Ground 신 필드 설계 문서 (승인됨)

Phase002 Step 1 의 산출물. Ground 확장 설계안과 그 근거.
사용자 리뷰 완료 (7개 포인트 승인 + 2개 조정).

---

## 0. 리뷰 조정 반영

| # | 원안 | 조정 | 근거 |
|---|------|------|------|
| 1 | Fullstack 확장 (iface + sqlc) | 승인 | — |
| 2 | 파서 위치 `pkg/parser/iface/`, `pkg/parser/sqlc/` | 승인 | — |
| 3 | Ground 필드 이름 | 승인 | — |
| 4 | 미사용 필드 제외 + `TableInfo.ColumnOrder` 추가 | 승인 | DDL 컬럼 순서 보존 필요 |
| 5 | populate 4개 분할 | 승인 | — |
| 6 | populate_symbol_table 공존 | **즉시 제거** | Phase003 으로 미루지 않고 본 Phase 에서 정리 |
| 7 | ErrStatus 출처 조사 | **완료** (아래 §7 참조) | opt2: populate_models 에서 FuncSpecs 조인 |

---

## 1. 조사 결론 요약

`internal/ssac/generator/` 약 250개 파일의 `SymbolTable` 사용 전수 조사 결과.

### 1.1 실제 사용되는 필드 (6개 중 4개)

| SymbolTable 필드 | Generator 읽기 | 이식 |
|-----------------|---------------|------|
| `Models` | **O** (Methods/Cardinality/Params/ErrStatus) | **O** |
| `DDLTables` | **O** (Columns 만) | **O** |
| `Operations` | **O** (일부 필드) | **O** |
| `RequestSchemas` | **O** (Fields/FieldConstraint) | **O** |
| `Funcs` | X (0건) | — |
| `DTOs` | X (0건) | — |

### 1.2 DDLTable 의 실제 소비 필드 (7개 중 1개)

`Columns` 만 사용. `ColumnOrder`, `ForeignKeys`, `Indexes`, `PrimaryKey`, `VarcharLen`, `CheckEnums` 전부 0건.

### 1.3 OperationSymbol 의 실제 소비

| 필드 | 읽기 |
|------|------|
| `PathParams` (Name, GoType) | **O** |
| `HasRequestBody` | **O** |
| `XPagination` (Style, DefaultLimit, MaxLimit) | **O** |
| `XSort` (Allowed, Default, Direction) | **O** |
| `XFilter` (Allowed) | **O** |
| `RequestFields`, `XInclude` | X |

### 1.4 MethodInfo 의 실제 소비

| 필드 | 읽기 |
|------|------|
| `Cardinality` ("one"/"many"/"exec") | **O** |
| `Params` | **O** |
| `ErrStatus` | **O** |
| `ParamTypes` | X |

### 1.5 접근 깊이

- 2단까지 (3단+ 없음).

### 1.6 비결정성 3곳

정렬 없는 map 순회:
- `lookup_call_err_status.go:21` — `for modelKey, ms := range st.Models`
- `lookup_ddl_type.go:9` — `for _, table := range st.DDLTables`
- `lookup_all_tables_column.go:8` — `for _, table := range st.DDLTables`

Ground 이식 시 **정렬 헬퍼** 제공.

### 1.7 원천 데이터 경로

| Ground 필드 | 원천 | 현재 Fullstack |
|------------|------|---------------|
| Models (from Go iface) | `specs/model/*.go` | **없음** (ModelDir 경로만) |
| Models (from sqlc queries) | `specs/db/queries/*.sql` | **없음** |
| Models (ErrStatus 주입) | FuncSpec `@error` | **있음** (`ProjectFuncSpecs`, `FullendPkgSpecs`) |
| Tables | `specs/db/*.sql` (DDL) | **있음** (`DDLTables`, `DDLResults`) |
| Ops | `specs/api/openapi.yaml` | **있음** (`OpenAPIDoc`) |
| ReqSchemas | `specs/api/openapi.yaml` | **있음** (`RequestConstraints`) |

**누락**: Go model 인터페이스 파싱, sqlc 쿼리 파싱 두 가지.

---

## 2. Fullstack 확장

### 2.1 신규 pkg/parser 패키지 2개

**`pkg/parser/iface/`** — Go 인터페이스 파싱 (`internal/ssac/validator/load_go_interfaces.go` 포트)

```go
// pkg/parser/iface/types.go
type Interface struct {
    Name    string   // "UserModel"
    Methods []Method // 순서 보존
}

type Method struct {
    Name        string
    Params      []Param
    ReturnType  string   // Go 반환 타입 문자열
    Cardinality string   // "one"/"many"/"exec" (시그니처 기반)
}

type Param struct {
    Name   string
    GoType string
}
```

**`pkg/parser/sqlc/`** — sqlc SQL 쿼리 파싱 (`internal/ssac/validator/load_sqlc_queries.go` + `extract_sqlc_params.go` 포트)

```go
// pkg/parser/sqlc/types.go
type Query struct {
    Name        string   // "-- name: UserFindByEmail :one" → "FindByEmail"
    Cardinality string   // ":one" / ":many" / ":exec"
    Model       string   // "User" (쿼리명 prefix 규약 유추)
    Params      []string // WHERE 절 파라미터 순서
    SQL         string   // 원본 SQL
}
```

**공개 API**:
- `iface.ParseDir(dir string) ([]Interface, []diagnostic.Diagnostic)`
- `sqlc.ParseDir(dir string) ([]Query, []diagnostic.Diagnostic)`

### 2.2 Fullstack 필드 추가

```go
// pkg/fullend/fullstack.go (확장)
type Fullstack struct {
    // ... 기존 필드 ...
    ModelInterfaces []iface.Interface   // 신규
    SqlcQueries     []sqlc.Query        // 신규
}
```

### 2.3 ParseAll 확장

`pkg/fullend/parse_all.go`:
- `ModelDir` 경로 존재 시 `iface.ParseDir(ModelDir)` → `fs.ModelInterfaces`
- `<specsDir>/db/queries/` 디렉토리 존재 시 `sqlc.ParseDir(...)` → `fs.SqlcQueries`

---

## 3. Ground 신 필드 설계

### 3.1 Ground 구조 확장

```go
// pkg/rule/ground.go (확장)
type Ground struct {
    // 기존 (보존)
    Lookup  map[string]StringSet
    Types   map[string]string
    Pairs   map[string]StringSet
    Config  map[string]bool
    Vars    StringSet
    Flags   StringSet
    Schemas map[string][]string

    // 신규 — generator 요구
    Models     map[string]ModelInfo
    Tables     map[string]TableInfo
    Ops        map[string]OperationInfo
    ReqSchemas map[string]RequestSchemaInfo
}
```

### 3.2 신 타입 정의 (`pkg/rule/ground_types.go`)

```go
// ModelInfo — Go 인터페이스 + sqlc + FuncSpec 결합
type ModelInfo struct {
    Name    string
    Methods map[string]MethodInfo
}

type MethodInfo struct {
    Cardinality string   // "one" / "many" / "exec"
    Params      []string // 매개변수 이름 순서
    ErrStatus   int      // HTTP 에러 코드 (@error 기반)
}

// TableInfo — DDL 테이블의 컬럼 + 순서
type TableInfo struct {
    Name        string
    Columns     map[string]string   // 컬럼명 → Go 타입
    ColumnOrder []string            // 원본 DDL 순서 (map 비결정성 해소)
}

// OperationInfo — OpenAPI operation 메타
type OperationInfo struct {
    ID             string
    Method         string            // "GET" 등 (참조용)
    Path           string            // 원본 path (참조용)
    PathParams     []PathParam
    HasRequestBody bool
    Pagination     *PaginationSpec   // nil 허용
    Sort           *SortSpec         // nil 허용
    Filter         *FilterSpec       // nil 허용
}

type PathParam struct {
    Name   string   // PascalCase (e.g. "CourseID")
    GoType string
}

type PaginationSpec struct {
    Style        string   // "offset" / "cursor"
    DefaultLimit int
    MaxLimit     int
}

type SortSpec struct {
    Allowed   []string
    Default   string
    Direction string
}

type FilterSpec struct {
    Allowed []string
}

// RequestSchemaInfo — OpenAPI requestBody 필드 제약
type RequestSchemaInfo struct {
    Fields map[string]FieldConstraint
}

type FieldConstraint struct {
    Required  bool
    Format    string
    MinLength *int
    MaxLength *int
    Minimum   *float64
    Maximum   *float64
    Pattern   string
    Enum      []string
}
```

### 3.3 미수용 필드 (Ground 에서 제외, YAGNI)

- `Funcs map[string]bool`, `DTOs map[string]bool` — generator 미사용
- `DDLTable.{ForeignKeys, Indexes, PrimaryKey, VarcharLen, CheckEnums}` — generator 미사용
- `OperationSymbol.{RequestFields, XInclude}` — generator 미사용
- `MethodInfo.ParamTypes` — generator 미사용

### 3.4 정렬 헬퍼 (비결정성 해소)

```go
// pkg/ground/sort_helpers.go (신설)
func SortedModelKeys(g *rule.Ground) []string { ... }
func SortedTableKeys(g *rule.Ground) []string { ... }
func SortedOpKeys(g *rule.Ground) []string { ... }
```

Generator (Phase004 에서 이식) 가 이 헬퍼를 통해 순회 → 산출물 결정적.

---

## 4. populate 함수 설계

### 4.1 신규 populate 4개

| 파일 | 입력 (Fullstack) | 출력 | 로직 |
|------|-----------------|------|------|
| `populate_models.go` | `ModelInterfaces`, `SqlcQueries`, `FullendPkgSpecs`, `ProjectFuncSpecs` | `g.Models` | iface → 초기 Methods / sqlc → 병합 / FuncSpec `@error` → ErrStatus 주입 |
| `populate_tables.go` | `DDLTables`, `DDLResults` | `g.Tables` | 컬럼 매핑 + **순서 보존** |
| `populate_ops.go` | `OpenAPIDoc` | `g.Ops` | path/query/security/x-pagination/x-sort/x-filter 수집 |
| `populate_request_schemas.go` | `RequestConstraints` (이미 파싱됨) | `g.ReqSchemas` | 필드 제약 구조 변환 |

### 4.2 build.go 갱신

```go
// pkg/ground/build.go (갱신)
func Build(fs *fullend.Fullstack) *rule.Ground {
    g := &rule.Ground{ ... }

    // 기존 populate
    populateOpenAPI(g, fs)
    populateSSaC(g, fs)
    populateStates(g, fs)
    populateFunc(g, fs)
    populateManifest(g, fs)
    populateDDL(g, fs)
    populateRego(g, fs)
    populateOpenAPIConstraints(g, fs)
    populateOpenAPIParams(g, fs)
    // populateSymbolTable 는 본 Phase 에서 삭제됨 (아래 §6 참조)
    populateVarTypes(g, fs)
    populateGoReservedWords(g)
    populateHurl(g, fs)

    // 신규
    populateModels(g, fs)
    populateTables(g, fs)
    populateOps(g, fs)
    populateRequestSchemas(g, fs)

    return g
}
```

---

## 5. populate_models 상세 (의사코드)

ErrStatus 조사(§7) 결과 반영한 3단계 병합 로직:

```go
// pkg/ground/populate_models.go
func populateModels(g *rule.Ground, fs *fullend.Fullstack) {
    g.Models = make(map[string]rule.ModelInfo)

    // 단계 1: Go interface → Models
    for _, iface := range fs.ModelInterfaces {
        methods := make(map[string]rule.MethodInfo, len(iface.Methods))
        for _, m := range iface.Methods {
            methods[m.Name] = rule.MethodInfo{
                Cardinality: m.Cardinality,
                Params:      paramNames(m.Params),
            }
        }
        g.Models[iface.Name] = rule.ModelInfo{Name: iface.Name, Methods: methods}
    }

    // 단계 2: sqlc query → Models (병합)
    for _, q := range fs.SqlcQueries {
        info, ok := g.Models[q.Model]
        if !ok {
            info = rule.ModelInfo{Name: q.Model, Methods: make(map[string]rule.MethodInfo)}
        }
        if info.Methods == nil {
            info.Methods = make(map[string]rule.MethodInfo)
        }
        info.Methods[q.Name] = rule.MethodInfo{
            Cardinality: stripColon(q.Cardinality), // ":one" → "one"
            Params:      q.Params,
        }
        g.Models[q.Model] = info
    }

    // 단계 3: FuncSpec @error → pkg._func 모델의 ErrStatus 주입
    // 참조: internal/orchestrator/inject_func_err_status_from_parsed.go 와 동형
    injectErrStatus(g, fs.FullendPkgSpecs)
    injectErrStatus(g, fs.ProjectFuncSpecs)
}

func injectErrStatus(g *rule.Ground, specs []funcspec.FuncSpec) {
    for _, spec := range specs {
        if spec.ErrStatus == 0 || spec.Package == "" {
            continue
        }
        modelKey := spec.Package + "._func"
        info, ok := g.Models[modelKey]
        if !ok {
            info = rule.ModelInfo{Name: modelKey, Methods: make(map[string]rule.MethodInfo)}
        }
        if info.Methods == nil {
            info.Methods = make(map[string]rule.MethodInfo)
        }
        funcName := upperFirst(spec.Name)
        mi := info.Methods[funcName]
        mi.ErrStatus = spec.ErrStatus
        info.Methods[funcName] = mi
        g.Models[modelKey] = info
    }
}
```

**핵심 규약**:
- 일반 모델 키: `"User"`, `"Course"` (iface/sqlc 유래)
- `@call` 함수 모델 키: `"auth._func"`, `"billing._func"` (FuncSpec 유래, 모든 메서드가 ErrStatus 보유)
- 함수명 대문자 시작 (`spec.Name[0]` → upper)

---

## 6. populate_symbol_table 즉시 제거 (원안에서 조정)

**기존 결정**: Phase002 에서 공존, Phase003 에서 제거.
**조정**: **본 Phase 에서 즉시 제거**.

### 조정 사유

- `populate_symbol_table.go` 는 `g.Lookup["SymbolTable.model"]` 에 모델명 집합만 넣음
- 실측: `pkg/validate/ssac`, `pkg/crosscheck`, `pkg/ground` 어디에서도 **사용처 0건**
- validate 에서 "이 모델 존재?" 조회는 신규 `g.Models[name]` 존재 여부로 동등 표현 가능

### 조정 실행 절차

**Sub-step D (build.go 통합)** 에 병행:

1. `pkg/ground/populate_symbol_table.go` 삭제
2. `pkg/ground/build.go` 에서 호출 라인 제거
3. `pkg/validate/**`, `pkg/crosscheck/**` 에서 `g.Lookup["SymbolTable.model"]` 참조 grep
4. 참조 0건 확인 → commit
5. 참조 있으면 → 해당 사용처를 `g.Models` 기반으로 치환 (소규모)

---

## 7. ErrStatus 출처 조사 결과 (승인 조정)

### 7.1 조사 결과

**체인**:

1. `pkg/parser/funcspec/apply_annotation.go:14-16` — `@error NNN` 어노테이션 파싱
   ```go
   case strings.HasPrefix(line, "@error "):
       if code, err := strconv.Atoi(...); err == nil {
           spec.ErrStatus = code
       }
   ```

2. `pkg/parser/funcspec/func_spec.go:10` — `FuncSpec.ErrStatus int` 필드

3. `internal/orchestrator/inject_func_err_status_from_parsed.go:15-35` — 주입
   ```go
   func injectFuncErrStatusFromParsed(st *ssacvalidator.SymbolTable, parsed *genapi.ParsedSSOTs) {
       allSpecs := append(parsed.FullendPkgSpecs, parsed.ProjectFuncSpecs...)
       for _, fs := range allSpecs {
           if fs.ErrStatus == 0 || fs.Package == "" { continue }
           modelKey := fs.Package + "._func"
           // ... ms.Methods[funcName].ErrStatus = fs.ErrStatus
           st.Models[modelKey] = ms
       }
   }
   ```

### 7.2 설계 결정

**옵션 2 채택**: `populate_models` 가 `fs.FullendPkgSpecs + fs.ProjectFuncSpecs` 를 입력으로 받아 ErrStatus 주입 (internal 로직과 동형).

대안 (옵션 1: iface 파서가 어노테이션 파싱) 은 기각 — iface 파서는 Go interface 정의만 다루는 게 책임 경계상 맞음. `@error` 는 FuncSpec 의 영역.

### 7.3 비고

`internal/ssac/validator/validate_err_status.go:15` 에서 `seq.ErrStatus` 를 읽는 코드는 **seq 레벨 검증용**. 모델 메서드 ErrStatus 와 무관 — populate_models 영역 밖.

---

## 8. 구현 순서 (Phase002 Step 2~5 세부화)

### 8.1 Part A — Fullstack 확장

**Sub-step A1.** `pkg/parser/iface/` 신설
- `internal/ssac/validator/load_go_interfaces.go` 포팅
- 공개: `ParseDir(dir)` → `([]Interface, []diagnostic.Diagnostic)`
- 단위 테스트 (dummy gigbridge 기반)

**Sub-step A2.** `pkg/parser/sqlc/` 신설
- `internal/ssac/validator/load_sqlc_queries.go` + `extract_sqlc_params.go` 포팅
- 공개: `ParseDir(dir)` → `([]Query, []diagnostic.Diagnostic)`
- 단위 테스트

**Sub-step A3.** Fullstack + ParseAll 배선
- `pkg/fullend/fullstack.go` 필드 추가
- `pkg/fullend/parse_all.go` 호출 추가

**커밋**: `feat(parser): iface + sqlc 파서 신설 + Fullstack 확장`

### 8.2 Part B — Ground 타입

- `pkg/rule/ground.go` 필드 4개 추가
- `pkg/rule/ground_types.go` 신설 (11개 타입)

**커밋**: `feat(rule): Ground generate 용 구조적 필드 타입 추가`

### 8.3 Part C — populate 4개 (독립 커밋)

- `feat(ground): populate_models 추가` (iface + sqlc + FuncSpec 병합)
- `feat(ground): populate_tables 추가` (DDL 컬럼 + 순서)
- `feat(ground): populate_ops 추가` (OpenAPI operation 메타)
- `feat(ground): populate_request_schemas 추가` (필드 제약)

각 커밋마다 단위 테스트 포함.

### 8.4 Part D — Build 통합 + populate_symbol_table 제거

- `pkg/ground/build.go` 업데이트 (신 4개 추가 + 기존 1개 제거)
- `pkg/ground/populate_symbol_table.go` 삭제
- 정렬 헬퍼 `pkg/ground/sort_helpers.go` 신설
- 기존 테스트 전수 통과 확인

**커밋**: `refactor(ground): 신 populate 통합 + populate_symbol_table 제거`

### 8.5 Part E — 검증

- `go build ./pkg/... ./internal/... ./cmd/...` 통과
- `go vet` 통과
- `go test ./pkg/...` 전부 통과
- dummy gigbridge 로드 후 신 필드 실값 검증 (assert 기반 단위 테스트)

---

## 9. 위험과 대응

### R1. sqlc 쿼리 → 모델 이름 매핑 규약

**질문**: `-- name: UserFindByEmail :one` 같은 쿼리명에서 "User" 와 "FindByEmail" 을 어떻게 분리?

**대응**: Sub-step A2 진입 시 `internal/ssac/validator/load_sqlc_queries.go` 의 로직 그대로 포팅. 새 규약 도입 안 함.

### R2. 파서 이식의 완전성

Go iface / sqlc 파서가 internal 과 동등한 결과를 내는지 **비교 테스트**:
- internal validator 가 만든 SymbolTable 의 Models 수집
- pkg iface+sqlc+populate_models 결과 비교
- 차이 있으면 파서 보강

### R3. Fullstack 확장의 파급

- `internal/ssac/parser`, `internal/ssac/validator` 는 건드리지 않음 → internal orchestrator 경로 영향 없음
- `internal/orchestrator/parsed.go` 가 아직 `internal/genapi.ParsedSSOTs` 사용 → 신 필드 소비 안 함. Phase004 에서 배선 교체 시 활용.

### R4. populate_symbol_table 제거의 안전성

- 실측 0건이지만 **삭제 직전 재확인**: `grep -rn "SymbolTable.model" pkg/ internal/`
- 발견 시 해당 위치를 `g.Models[name]` 기반으로 먼저 치환

---

## 10. 산출물 체크리스트 (Definition of Done)

- [ ] `pkg/parser/iface/` 패키지 (types, parse_dir, 테스트)
- [ ] `pkg/parser/sqlc/` 패키지 (types, parse_dir, extract_params, 테스트)
- [ ] `pkg/fullend/fullstack.go` 에 `ModelInterfaces`, `SqlcQueries` 필드
- [ ] `pkg/fullend/parse_all.go` 가 두 파서 호출
- [ ] `pkg/rule/ground.go` 에 `Models`, `Tables`, `Ops`, `ReqSchemas` 필드
- [ ] `pkg/rule/ground_types.go` 에 신 타입 11개
- [ ] `pkg/ground/populate_models.go` (3단계 병합)
- [ ] `pkg/ground/populate_tables.go` (컬럼 + ColumnOrder)
- [ ] `pkg/ground/populate_ops.go`
- [ ] `pkg/ground/populate_request_schemas.go`
- [ ] `pkg/ground/sort_helpers.go` (SortedModelKeys 등)
- [ ] `pkg/ground/build.go` 에 신 populate 통합
- [ ] `pkg/ground/populate_symbol_table.go` **삭제**
- [ ] 단위 테스트 (파서 2개 + populate 4개)
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] `go test ./pkg/...` 전부 통과 (특히 기존 validate/crosscheck 테스트 무영향)
- [ ] dummy gigbridge 로드 시 신 필드 비어있지 않음

---

## 11. 다음 Phase 예고

- **Phase003** — validate/crosscheck 에서 **신 필드로 점진 마이그** (원안의 평탄 populate 제거 작업에서 SymbolTable 부분은 본 Phase 에서 이미 완료됨).
- **Phase004** — generator 이식 시 Ground 신 필드 직접 소비.
- **Phase005** — dummy 회귀 검증.
