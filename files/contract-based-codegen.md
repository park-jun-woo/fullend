# Contract-Based Code Generation

## 문제

현재 codegen은 **all-or-nothing**이다.

- `fullend gen` → artifacts/ 전체 덮어씀
- 생성된 코드를 수동 수정하면 다음 gen에서 소실
- 코드젠 품질이 90%면 나머지 10%를 고칠 방법이 없음

개발자의 선택지: ① 코드젠 자체를 고치거나 ② 생성된 코드를 포크하고 다시는 gen을 못 돌리거나.

## 핵심 아이디어

**함수의 입출력 계약(contract)만 유지하면 내부 구현(body)은 수정해도 된다.**

```
SSOT (specs/)          계약 (contract)            구현 (body)
──────────────        ──────────────────         ──────────────
fullend 소유           fullend 소유               개발자 소유
수정 → gen 재실행      SSOT에서 파생, 자동 갱신    수동 수정 보존
```

## 계약의 정의

계약(contract)은 **Go 함수 시그니처가 아니라 SSOT에서 파생된 입출력 명세**다.

gin 핸들러의 Go 시그니처는 항상 `func (h *Handler) XXX(c *gin.Context)`로 동일하므로, Go 시그니처를 해싱하면 의미가 없다. 대신 SSOT가 정의하는 실제 계약을 해싱한다.

### Service Handler

```
계약 출처: SSaC 시퀀스 + OpenAPI schema
해싱 대상: operationId + SSaC @시퀀스 타입 목록 + request fields + response fields

예: CreateGig
    hash("CreateGig|@post,@response|title:string,budget:int64|gig:Gig") → a3f8c1
```

body (수정 가능): 요청 파싱, 모델 호출, 에러 처리, 응답 조립.

### Model Implementation

```
계약 출처: models_gen.go interface 시그니처 (SSaC + DDL에서 파생)
해싱 대상: 함수명 + 파라미터 타입 목록 + 반환 타입

예: gigModelImpl.Create
    hash("Create|string,int64,int64|*Gig,error") → e1d9f2
```

body (수정 가능): SQL 쿼리, 트랜잭션, 캐시 로직.

### State Machine

```
계약 출처: Mermaid stateDiagram
해싱 대상: state 목록 + transition 목록

예: gig state machine
    hash("draft,open,in_progress,completed|PublishGig:draft→open,...") → f5b3a9
```

body (수정 가능): CanTransition 내부 검증 로직.

### Middleware

```
계약 출처: fullend.yaml claims config
해싱 대상: CurrentUser struct fields

예: BearerAuth
    hash("ID:int64,Email:string,Role:string") → c2d4e6
```

body (수정 가능): 토큰 추출, 검증, claims 매핑 로직.

## 소유권 디렉티브: `//fullend:`

생성된 Go 코드 자체에 메타 정보를 내장한다. 외부 lock 파일 불필요.

### 디렉티브 형식

```go
//fullend:<ownership> ssot=<path> contract=<hash>
```

| 필드 | 값 | 의미 |
|---|---|---|
| ownership | `gen` | fullend 소유 — gen 시 덮어씀 |
| ownership | `preserve` | 개발자 소유 — gen 시 body 보존 |
| `ssot=` | SSOT 파일 상대경로 | 이 함수의 출처 |
| `contract=` | SSOT 파생 계약의 SHA256 (7자리) | 계약 변경 감지용 해시 |

### 부착 위치

**함수** — 함수 doc comment 위치:

```go
//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1
func (h *Handler) CreateGig(c *gin.Context) { ... }
```

**파일 레벨** (const, var, type 블록) — 파일 첫 줄 코멘트:

```go
//fullend:gen ssot=states/gig.md contract=f5b3a9
package gigstate

const (
    StateDraft      = "draft"
    StateOpen       = "open"
)

var transitions = map[transitionKey]string{ ... }
```

파일 레벨 디렉티브는 파일 전체를 하나의 단위로 관리한다. 함수 레벨 디렉티브가 있으면 함수 레벨이 우선.

### 생성 직후 (fullend 소유)

```go
//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1
func (h *Handler) CreateGig(c *gin.Context) {
    var req struct {
        Title  string `json:"title"`
        Budget int64  `json:"budget"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    gig, err := h.GigModel.Create(req.Title, req.Budget)
    if err != nil {
        c.JSON(500, gin.H{"error": "create failed"})
        return
    }
    c.JSON(200, gin.H{"gig": gig})
}
```

### 개발자가 수정 후 (개발자 소유)

`gen` → `preserve`로 바꾸는 것만으로 소유권 전환:

```go
//fullend:preserve ssot=service/gig/create_gig.ssac contract=a3f8c1
func (h *Handler) CreateGig(c *gin.Context) {
    var req struct {
        Title  string `json:"title"`
        Budget int64  `json:"budget"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 수동 추가: 중복 제목 검사
    if exists, _ := h.GigModel.ExistsByTitle(req.Title); exists {
        c.JSON(409, gin.H{"error": "duplicate title"})
        return
    }
    gig, err := h.GigModel.Create(req.Title, req.Budget)
    if err != nil {
        c.JSON(500, gin.H{"error": "create failed"})
        return
    }
    c.JSON(200, gin.H{"gig": gig})
}
```

## 한 파일에 gen/preserve 혼재 처리

`model/gig.go`처럼 한 파일에 여러 함수가 있을 때, 일부만 preserve할 수 있다:

```go
package model

//fullend:preserve ssot=db/gigs.sql contract=e1d9f2
func (m *gigModelImpl) Create(title string, budget int64, categoryID int64) (*Gig, error) {
    // 개발자가 수정한 SQL — 보존됨
}

//fullend:gen ssot=db/gigs.sql contract=b8c3a7
func (m *gigModelImpl) FindByID(id int64) (*Gig, error) {
    // fullend가 재생성 — 덮어씀
}
```

### 구현: Go AST 함수 단위 splice

파일 전체를 덮어쓰는 대신, **함수 단위로 교체/보존**한다:

1. 기존 파일을 Go AST로 파싱 → 함수별 디렉티브 + body 추출
2. 새로 생성할 코드를 Go AST로 파싱 → 함수별 디렉티브 + body 추출
3. 함수별로 gen/preserve 판정
4. import 문은 사용된 함수들의 합집합으로 재계산
5. `go/format`으로 출력

## fullend gen 동작

### 기본 동작 (preserve 모드)

```bash
fullend gen specs/ artifacts/
```

각 함수에 대해:

| 디렉티브 | contract 변경 | gen 동작 |
|---|---|---|
| 없음 (새 파일) | — | 생성 + `//fullend:gen` 부착 |
| `gen` | — | 덮어씀 |
| `preserve` | 없음 | **스킵 (body 보존)** |
| `preserve` | 있음 | **충돌 경고** |

### 전체 초기화 (--reset)

```bash
fullend gen --reset specs/ artifacts/

⚠ --reset: 모든 preserve 함수가 초기화됩니다.
  preserve 함수 5개가 fullend:gen으로 되돌아갑니다.
  계속하시겠습니까? (Y/n): _
```

`Y` 또는 Enter 시 모든 `//fullend:preserve`를 `//fullend:gen`으로 바꾸고 body 재생성. `n` 입력 시 중단.

### 충돌 시 출력

```
⚠ Contract changed: service/gig/create_gig.go:Handler.CreateGig
  SSOT:     service/gig/create_gig.ssac (modified)
  Contract: a3f8c1 → b7d2e5 (arg 'deadline' added)
  Action:   body preserved, manual merge required
            new contract written to create_gig.go.new
```

개발자는 `.new` 파일을 참고해 수동 머지 후, 원본의 `contract=` 해시를 새 값으로 갱신.

## fullend contract 커맨드

```bash
fullend contract specs/ artifacts/

Contract Status:
  gen       service/gig/list_gigs.go      ListGigs         fullend 소유
  gen       model/gig.go                  FindByID         fullend 소유
  gen       states/gigstate/gigstate.go   CanTransition    fullend 소유
  preserve  service/gig/create_gig.go     CreateGig        계약 유지 ✓
  preserve  model/gig.go                  Create           계약 유지 ✓
  preserve  service/gig/update_gig.go     UpdateGig        계약 위반 ✗ (arg added)
  orphan    service/gig/old_feature.go    OldFeature       SSOT 삭제됨 ⚠

  gen:      3 functions (fullend 소유)
  preserve: 2 functions (개발자 소유, 계약 유지)
  broken:   1 function (계약 위반, 수동 머지 필요)
  orphan:   1 function (출처 SSOT 없음)
```

## validate 통합

`fullend validate`에서 `//fullend:` 디렉티브를 파싱하여 계약 위반을 검출:

```
✓ Config       my-project
✓ OpenAPI      7 endpoints
✓ DDL          3 tables
✓ SSaC         7 functions
✓ Cross        0 mismatches
✗ Contract     1 violation, 1 orphan
               service/gig/update_gig.go:UpdateGig — contract a3f8c1 → b7d2e5
               service/gig/old_feature.go:OldFeature — ssot not found
```

## Feature Chain 통합

```bash
fullend chain CreateGig specs/

── Feature Chain: CreateGig ──

  OpenAPI    api/openapi.yaml:45         POST /gigs
  SSaC       service/gig/create_gig.ssac @post @response
  DDL        db/gigs.sql:1               CREATE TABLE gigs
  Rego       policy/authz.rego:12        resource: gig
  Gherkin    scenario/gig.feature:5      Scenario: Create a new gig

  ── Artifacts ──
  Handler    internal/service/gig/create_gig.go:CreateGig    preserve ✎
  Model      internal/model/gig.go:Create                    preserve ✎
  Authz      internal/authz/authz.go                         gen
  Types      internal/model/types.go:Gig                     gen
```

SSOT 노드뿐 아니라 파생된 artifacts 함수까지 chain에 포함. 소유권 상태도 표시.

## 프론트엔드 (TSX)

TSX 파일에서는 모듈 최상위(JSX 외부)이므로 JS 코멘트를 사용:

```tsx
// fullend:gen ssot=frontend/gig_list.html contract=d4e5f6
export function GigListPage() {
    // 생성된 React 컴포넌트
}
```

```tsx
// fullend:preserve ssot=frontend/gig_list.html contract=d4e5f6
export function GigListPage() {
    // 개발자가 수정한 컴포넌트
}
```

Go와 동일한 `// fullend:` 접두사. 파싱: 정규식으로 `// fullend:` 패턴 추출.

## 계약 감지 구현

### 디렉티브 파싱

```go
// Go AST + 코멘트에서 //fullend: 디렉티브 추출
func parseDirective(fn *ast.FuncDecl) *Directive {
    if fn.Doc == nil { return nil }
    for _, c := range fn.Doc.List {
        if strings.HasPrefix(c.Text, "//fullend:") {
            return parseFullendComment(c.Text)
        }
    }
    return nil
}

type Directive struct {
    Ownership string // "gen" or "preserve"
    SSOT      string // relative path to SSOT file
    Contract  string // 7-char SHA256 of SSOT-derived contract
}
```

### contract hash 계산

```go
// SSOT에서 계약 해싱 — Go 시그니처가 아닌 SSOT 명세를 해싱
func computeContractHash(sf ssacparser.ServiceFunc) string {
    // SSaC 시퀀스 타입 + request fields + response fields
    var parts []string
    parts = append(parts, sf.OperationID)
    for _, seq := range sf.Sequences {
        parts = append(parts, seq.Type)  // @get, @post, @auth, ...
    }
    for _, arg := range sf.Args {
        parts = append(parts, arg.Name+":"+arg.Type)
    }
    // response fields ...
    h := sha256.Sum256([]byte(strings.Join(parts, "|")))
    return hex.EncodeToString(h[:])[:7]
}
```

### gen 시 함수 단위 splice

```go
func generateWithPreserve(oldPath, newContent string) (string, []Warning) {
    oldFuncs := parseFuncsWithDirective(oldPath)
    newFuncs := parseFuncsWithDirective(newContent)
    var warnings []Warning

    for name, newFn := range newFuncs {
        oldFn, ok := oldFuncs[name]
        if !ok { continue }  // 새 함수 → 그대로 생성

        if oldFn.Directive.Ownership == "preserve" {
            if oldFn.Directive.Contract == newFn.Directive.Contract {
                // 계약 유지 → body 보존, 디렉티브 유지
                newFn.Body = oldFn.Body
                newFn.Directive.Ownership = "preserve"
            } else {
                // 계약 변경 → body 보존 + 경고 + .new 파일
                warnings = append(warnings, Warning{
                    File:        oldPath,
                    Function:    name,
                    OldContract: oldFn.Directive.Contract,
                    NewContract: newFn.Directive.Contract,
                })
                newFn.Body = oldFn.Body
                newFn.Directive.Ownership = "preserve"
            }
        }
        // ownership == "gen" → newFn 그대로 (덮어씀)
    }

    // import 재계산: 보존된 body + 새 body에서 사용하는 패키지 합집합
    imports := collectImportsFromFuncs(newFuncs)
    return renderFile(imports, newFuncs), warnings
}
```

## 이것이 해결하는 문제

| 현재 | Contract-Based |
|---|---|
| gen 결과 90% → 나머지 10% 수정 불가 | 90% gen + 10% 수동 수정, 공존 |
| SSOT 변경 → gen 재실행 → 수동 수정 소실 | 계약 유지 시 body 보존 |
| 코드젠 버그 → 도구 자체를 고쳐야 함 | 버그 있는 body만 수동 수정으로 우회 |
| artifacts를 git 추적할 이유 없음 | artifacts를 git 추적하는 의미 생김 |
| Feature Chain = SSOT만 | Feature Chain = SSOT + artifacts 통합 |
