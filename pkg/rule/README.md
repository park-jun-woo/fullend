# pkg/rule

fullend validate/crosscheck 공통 toulmin 룰 함수.

모든 함수는 `func(claim, ground, backing any) (bool, any)` 시그니처를 따른다.
backing이 판정 기준을 담고, 같은 함수를 backing만 바꿔 재사용한다.

## 공통 타입

### Ground

모든 룰 함수가 공유하는 검증 컨텍스트. 호출자가 그래프 평가 전에 구성한다.

```go
type Ground struct {
    Lookup  map[string]StringSet // "target.kind" -> set of names
    Types   map[string]string    // "target.kind.name" -> type string
    Pairs   map[string]StringSet // "target.pairKind" -> set of "key:value"
    Config  map[string]bool      // config key -> present
    Vars    StringSet            // declared variable names
    Flags   StringSet            // flags for defeaters
    Schemas map[string][]string  // "target.schema" -> ordered field list
}
```

**LookupKey 컨벤션**: `"SSOT.kind"` 또는 `"SSOT.kind.scope"` 형식.

```
"OpenAPI.operationId"          → {"CreateGig", "ListGigs"}
"DDL.column.users"             → {"id", "email", "name"}
"SSaC.funcName"                → {"CreateGig", "ListGigs"}
"Config.claims"                → {"ID", "Email", "Role"}
"Func.request.HoldEscrow"     → {"GigID", "Amount", "ClientID"}
```

### Evidence

위반 시 반환되는 결과 상세.

```go
type Evidence struct {
    Rule    string // 규칙 ID: "X-1", "S-48" 등
    Level   string // "ERROR" 또는 "WARNING"
    Ref     string // 검사 대상 이름
    Message string // 에러 메시지 템플릿
}

type SchemaEvidence struct {
    Rule    string
    Level   string
    Missing []string // 누락된 필드 목록
    Message string
}
```

## Warrant 함수

### RefExists

A가 참조하는 이름이 B에 존재해야 한다.

- **claim**: `string` (참조 이름)
- **ground**: `*Ground` (Lookup 사용)
- **evidence**: `*Evidence`

```go
type RefExistsBacking struct {
    LookupKey string // Ground.Lookup 키: "OpenAPI.operationId", "DDL.column.users" 등
    Rule      string
    Level     string
    Message   string
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `OpenAPI.operationId` | X-15 | SSaC funcName → OpenAPI |
| `DDL.column.users` | X-1 | x-sort column → DDL |
| `SSaC.funcName` | X-23 | States transition → SSaC |
| `States.diagram` | X-24 | @state → diagram 존재 |
| `DDL.table` | X-31 | @ownership table → DDL |
| `OpenAPI.path` | X-35 | Hurl path → OpenAPI |
| `Func.spec` | X-39 | @call → func 구현 존재 |
| `Config.claims` | X-49, X-53 | currentUser/Rego → claims |
| `Config.roles` | X-63 | Rego role → Config |
| `OpenAPI.operationId` | TM-1 | STML data-fetch → OpenAPI |
| `SymbolTable.model` | S-48 | SSaC Model → 심볼 테이블 |

커버 규칙: ~31개.

### CoverageCheck

B에 정의된 것이 A에서 사용되어야 한다. RefExists의 역방향.

- **claim**: `string` (정의된 항목 이름)
- **ground**: `*Ground` (Lookup 사용)
- **evidence**: `*Evidence`

```go
type CoverageCheckBacking struct {
    LookupKey string // Ground.Lookup 키: 사용처의 이름 집합
    Rule      string
    Level     string
    Message   string
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `SSaC.funcName` | X-16 | OpenAPI operationId가 SSaC에서 사용되는지 |
| `SSaC.response` | X-18 | OpenAPI response field가 @response에 있는지 |
| `SSaC.modelRef` | X-55 | DDL table이 SSaC에서 참조되는지 |
| `SSaC.callRef` | X-62 | func spec이 @call에서 사용되는지 |
| `Rego.claims` | X-54 | Config claims가 Rego에서 참조되는지 |
| `Rego.roles` | X-64 | Config roles가 Rego에서 사용되는지 |

커버 규칙: ~9개.

### PairMatch

A의 (key:value) 쌍이 B에 매칭되어야 한다.

- **claim**: `string` ("key:value" 형식)
- **ground**: `*Ground` (Pairs 사용)
- **evidence**: `*Evidence`

```go
type PairMatchBacking struct {
    LookupKey string // Ground.Pairs 키
    Rule      string
    Level     string
    Message   string
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `Policy.auth` | X-28 | SSaC @auth (action:resource) → Rego allow |
| `SSaC.auth` | X-29 | Rego allow → SSaC @auth |
| `Config.middleware` | X-50 | OpenAPI security → Config middleware |
| `OpenAPI.security` | X-51 | Config middleware → OpenAPI security |
| `SSaC.subscribe` | X-57 | @publish topic → @subscribe |
| `SSaC.publish` | X-58 | @subscribe topic → @publish |

커버 규칙: ~6개.

### TypeMatch

A의 타입이 B의 타입과 일치해야 한다. 호출자가 타입을 정규화한 뒤 Ground.Types에 저장한다.

- **claim**: `*TypeClaim` (Name + SourceType)
- **ground**: `*Ground` (Types 사용)
- **evidence**: `*Evidence`

```go
type TypeMatchBacking struct {
    LookupKey string // Ground.Types 키 접두사 (+ ".name" 자동 연결)
    Rule      string
    Level     string
    Message   string
}

type TypeClaim struct {
    Name       string // 필드/파라미터 이름
    SourceType string // 소스 측 타입 (정규화 후)
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `Func.request.HoldEscrow` | X-44 | @call input type ↔ Request field type |
| `DDL.check.users` | X-69 | DDL CHECK values ↔ OpenAPI enum |

커버 규칙: ~4개.

### SchemaMatch

A의 필드 집합이 B의 스키마에 존재해야 한다.

- **claim**: `[]string` (소스 필드 이름 목록)
- **ground**: `*Ground` (Schemas 사용)
- **evidence**: `*SchemaEvidence` (Missing 포함)

```go
type SchemaMatchBacking struct {
    LookupKey string // Ground.Schemas 키
    Rule      string
    Level     string
    Message   string
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `OpenAPI.response.CreateGig` | X-17 | @response fields → OpenAPI response |
| `Func.request.HoldEscrow` | X-42 | @call Inputs → FuncRequest fields |
| `SSaC.publish.topic` | X-59 | @subscribe fields → @publish payload |
| `OpenAPI.required.CreateGig` | X-66 | SSaC used fields → OpenAPI required |

커버 규칙: ~7개.

### ConfigRequired

기능을 사용하면 Config에 해당 설정이 있어야 한다.

- **claim**: 무시 (nil 가능)
- **ground**: `*Ground` (Config 사용)
- **evidence**: `*Evidence`

```go
type ConfigRequiredBacking struct {
    ConfigKey string // Ground.Config 키: "backend.auth.claims", "queue.backend" 등
    Rule      string
    Level     string
    Message   string
}
```

| ConfigKey | 원래 규칙 | 설명 |
|---|---|---|
| `backend.auth.claims` | X-48 | currentUser 사용 시 필수 |
| `backend.middleware` | X-52 | endpoint security 사용 시 필수 |
| `queue.backend` | X-56 | @publish/@subscribe 사용 시 필수 |

커버 규칙: ~3개.

### FieldRequired

필드가 있어야 하거나(Required=true) 없어야 한다(Required=false). SSaC 전용.

- **claim**: `map[string]bool` (필드명 → 값 존재 여부)
- **ground**: 무시
- **evidence**: `*Evidence`

```go
type FieldRequiredBacking struct {
    SeqType  string // "@get", "@post", "@put", "@delete", "@empty", "@state", "@auth", "@call", "@publish"
    Field    string // "Model", "Result", "Inputs", "Target", "Message" 등
    Required bool   // true = 있어야 함, false = 없어야 함
    Rule     string
    Level    string
    Message  string
}
```

backing 하나당 필드 하나. 그래프에서 규칙 ID가 분리된다.

| SeqType | Field | Required | 원래 규칙 |
|---|---|---|---|
| @get | Model | true | S-1 |
| @get | Result | true | S-2 |
| @post | Model | true | S-3 |
| @post | Result | true | S-4 |
| @post | Inputs | true | S-5 |
| @put | Model | true | S-6 |
| @put | Result | **false** | S-7 |
| @put | Inputs | true | S-8 |
| @delete | Model | true | S-9 |
| @delete | Result | **false** | S-10 |
| @empty | Target | true | S-12 |
| @empty | Message | true | S-13 |
| @state | DiagramID | true | S-14 |
| @state | Inputs | true | S-15 |
| @state | Transition | true | S-16 |
| @state | Message | true | S-17 |
| @auth | Action | true | S-18 |
| @auth | Resource | true | S-19 |
| @auth | Message | true | S-20 |
| @call | Model | true | S-21 |
| @publish | Topic | true | S-23 |
| @publish | Payload | true | S-24 |

커버 규칙: ~22개.

### VarDeclared

변수는 선언 후 사용해야 한다. SSaC 전용.

- **claim**: `string` (변수명)
- **ground**: `*Ground` (Vars 사용)
- **evidence**: `*Evidence`

```go
type VarDeclaredBacking struct {
    Rule    string
    Level   string
    Message string
}
```

커버 규칙: ~4개 (S-27~S-30).

### ForbiddenRef

이름이 금지 목록에 있으면 안 된다. RefExists의 역논리.

- **claim**: `string` (검사 대상 이름)
- **ground**: `*Ground` (Lookup 사용 — 금지 이름 집합)
- **evidence**: `*Evidence`

```go
type ForbiddenRefBacking struct {
    LookupKey string // Ground.Lookup 키 (금지 이름 집합)
    Rule      string
    Level     string
    Message   string
}
```

| LookupKey 예시 | 원래 규칙 | 설명 |
|---|---|---|
| `go.reserved` | S-34, S-35 | Go 예약어 (type, select, func 등) |
| `ssac.reservedSource` | S-33 | 예약 소스 (config, currentUser, request, query) |
| `ssac.configPrefix` | S-31 | config.* 입력 금지 |
| `publish.forbidden` | S-32 | @publish에서 query 금지 |
| `subscribe.forbidden` | S-42, S-43 | @subscribe에서 request, query 금지 |
| `http.forbidden` | S-44 | HTTP 함수에서 message 금지 |

커버 규칙: ~7개.

### NameFormat

이름이 형식 규칙을 만족해야 한다.

- **claim**: `string` (검사 대상 이름)
- **ground**: 무시
- **evidence**: `*Evidence`

```go
type NameFormatBacking struct {
    Pattern string // "uppercase-start", "no-dot-prefix", "dot-method"
    Rule    string
    Level   string
    Message string
}
```

| Pattern | 원래 규칙 | 설명 |
|---|---|---|
| `uppercase-start` | S-46 | Result 타입 대문자 시작 |
| `no-dot-prefix` | S-47 | package-prefix @model 금지 |
| `dot-method` | S-26 | Model.Method 형식 필수 |

커버 규칙: ~3개.

## Defeater 함수

예외 조건이 충족되면 warrant를 무력화한다.
모두 `Ground.Flags` 기반. 호출자가 평가 전에 해당 플래그를 설정한다.

| 함수 | backing | Flags 키 | 면제 대상 |
|---|---|---|---|
| `IsSkipped` | `string` (SSOT kind) | `skipped.<kind>` | 해당 SSOT 관련 warrant 전체 |
| `IsPkgModel` | nil | `pkgModel` | RefExists(SSaC→DDL), CoverageCheck(DDL→SSaC) |
| `IsDTO` | nil | `dto` | RefExists(Model→DDL table) |
| `IsArchived` | nil | `archived` | CoverageCheck(DDL→SSaC) |
| `IsSensitiveCol` | nil | `sensitive` | CoverageCheck(DDL→OpenAPI) |
| `IsNoSensitive` | nil | `nosensitive` | RefExists(sensitive pattern) |
| `IsSubscribe` | nil | `subscribe` | FieldRequired(HTTP 전용), RefExists(request→OpenAPI) |
| `IsImplicitVar` | nil | `implicit.<name>` | VarDeclared (claim=변수명) |
| `IsCustomTS` | nil | `customTS.<name>` | RefExists(STML bind→OpenAPI) (claim=필드명) |

## 커버리지

| 함수 | 규칙 수 | 비율 |
|---|---|---|
| RefExists | ~31 | 22% |
| FieldRequired | ~22 | 15% |
| CoverageCheck | ~9 | 6% |
| ForbiddenRef | ~7 | 5% |
| SchemaMatch | ~7 | 5% |
| PairMatch | ~6 | 4% |
| TypeMatch | ~4 | 3% |
| VarDeclared | ~4 | 3% |
| ConfigRequired | ~3 | 2% |
| NameFormat | ~3 | 2% |
| **공통 합계** | **~96** | **68%** |
| 고유 규칙 | ~46 | 32% |
| Defeater | 9 | — |
