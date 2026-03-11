# Phase 021: 외부 서비스 모델 지원 + Func 순수성 강제

## 배경

현재 `@model`은 DDL 테이블(DB)에 1:1 종속. 외부 서비스 통신은 `@call func`으로만 가능.
`@call func`을 순수 로직(I/O 금지)으로 제한하기로 결정했으므로, 외부 서비스 I/O를 담당할 계층이 필요.

### @model vs @call 역할 구분 (확정)

- **@model** = I/O 추상화 (CRUD). DB 테이블, 외부 API, 파일 스토리지 등 모든 데이터 접근.
- **@call func** = 순수 비즈니스 로직. I/O 금지. 계산, 판단, 변환만 허용.
- 서비스 시퀀스가 @call(계산) → @model(저장) 순서를 조율.

## 현재 모델 흐름

```
DDL SQL → 테이블 정의
  ↓
SSaC @model Gig.FindByID(...) → 인터페이스 메서드
  ↓
fullend gluegen → *sql.DB 기반 구현체 자동 생성
```

## 목표

### A. 모델을 두 종류로 확장

| 종류 | 데이터 소스 | 정의 SSOT | 코드젠 결과 |
|---|---|---|---|
| DB 모델 | PostgreSQL 등 | DDL SQL | `*sql.DB` 기반 구현체 (기존) |
| 외부 모델 | 외부 API | 외부 서비스 스펙 (신규) | `*http.Client` 기반 구현체 (신규) |

### B. Func 순수성 강제 + TODO 레벨 수정

Phase 020에서 이관된 항목:

| # | 카테고리 | 문제 |
|---|---|---|
| 1 | Func TODO 레벨 | 본체 미구현이 WARNING이지만 ERROR가 맞음 |
| 2 | Func 순수성 검증 | @call func이 I/O 패키지를 import하면 ERROR |

외부 모델이 제공되어야 Func 순수성 강제가 의미 있음 (I/O를 @model로 대체할 수 있으니까).

## 설계

### 1. 외부 서비스 스펙 SSOT (신규)

`specs/<project>/external/` 디렉토리에 서비스 계약 정의.
Go 인터페이스 파일 + 어노테이션 형식:

```go
// specs/dummy-gigbridge/external/escrow.go
package external

// @service escrow
// @base_url ${ESCROW_SERVICE_URL}

// @endpoint POST /api/escrow/hold
type HoldRequest struct {
    GigID    int64 `json:"gig_id"`
    Amount   int64 `json:"amount"`
    ClientID int64 `json:"client_id"`
}

type HoldResponse struct {
    TransactionID int64  `json:"transaction_id"`
    Status        string `json:"status"`
}

// @endpoint POST /api/escrow/release
type ReleaseRequest struct {
    GigID        int64 `json:"gig_id"`
    Amount       int64 `json:"amount"`
    FreelancerID int64 `json:"freelancer_id"`
}

type ReleaseResponse struct {
    TransactionID int64  `json:"transaction_id"`
    Status        string `json:"status"`
}
```

### 2. SSaC에서의 사용

SSaC 서비스 코드에서는 DB 모델과 동일한 `@model` 문법:

```go
// @post HoldResponse result = Escrow.Hold({GigID: gig.ID, Amount: gig.Budget, ClientID: currentUser.ID})
```

SSaC는 `Escrow`가 DB인지 외부인지 **모름** — fullend가 판단.

### 3. fullend의 판단 로직

모델명이 주어졌을 때 fullend가 판별하는 우선순위:

1. DDL 테이블에 매핑되는가? (`escrows` 테이블 존재?) → DB 모델
2. 외부 서비스 스펙에 정의되어 있는가? (`external/escrow.go`에 `@service escrow`?) → 외부 모델
3. 어디에도 없으면 → ERROR

### 4. 코드젠 분기

```
모델 수집 → DB 모델 / 외부 모델 분류
  ├── DB 모델 → 기존 model_impl.go 로직 (SQL 기반)
  └── 외부 모델 → 신규 로직 (HTTP client 기반)
```

**DB 모델 (기존):**
```go
type GigModel interface {
    FindByID(id int64) (*Gig, error)
}
type gigModelImpl struct { db *sql.DB }
```

**외부 모델 (신규):**
```go
type EscrowModel interface {
    Hold(req EscrowHoldRequest) (*EscrowHoldResponse, error)
}
type escrowModelImpl struct {
    baseURL string
    client  *http.Client
}
func (m *escrowModelImpl) Hold(req EscrowHoldRequest) (*EscrowHoldResponse, error) {
    // POST ${ESCROW_SERVICE_URL}/api/escrow/hold
    // JSON marshal req → HTTP request → JSON unmarshal response
}
```

### 5. Func 순수성 강제

`@call func`은 순수 비즈니스 로직만 허용. DB, 네트워크, 파일 I/O 접근 금지.

funcspec 파서가 이미 Go AST로 파싱하므로 import 목록 추출 가능. `FuncSpec`에 `Imports []string` 필드를 추가하고, crosscheck에서 금지 import를 검사.

**금지 import (ERROR):**
- DB: `database/sql`, `github.com/lib/pq`, `github.com/jackc/pgx`
- 네트워크: `net/http`, `net/rpc`, `google.golang.org/grpc`
- 파일 I/O: `os`, `io`, `io/ioutil`, `bufio`

**허용:**
- `math`, `strings`, `strconv`, `fmt`, `time`, `crypto/*`, `encoding/*`, `regexp`, `sort`, `errors` 등 순수 패키지

### 6. Func TODO 레벨 수정

본체 미구현(`// TODO: implement`만 있는 stub)이 `WARNING`으로 분류됨. 실행 시 빈 값을 반환하므로 `ERROR`가 맞음.

**수정**: func.go:83 `Level: "WARNING"` → `Level: "ERROR"`

### 7. 교차 검증 변경

| 규칙 | 현재 | 변경 |
|---|---|---|
| @model → DDL | 모든 @model에 DDL 테이블 요구 | DDL 또는 외부 스펙 중 하나 존재하면 OK |
| @result ↔ DDL | 외부 타입은 패키지 접두사로 스킵 | 외부 모델 응답 타입은 외부 스펙에서 검증 |
| DDL → SSaC | DDL 테이블 미참조 WARN | 변경 없음 |
| External → SSaC | 없음 (신규) | 외부 스펙 정의 후 SSaC 미참조 시 WARN |
| Func 순수성 | 없음 (신규) | 금지 import 감지 시 ERROR |
| Func TODO | WARNING | ERROR |

### 8. main.go 의존성 주입

```go
// DB 모델
server.Gig = &gigsvc.Handler{
    GigModel: model.NewGigModel(conn),
}

// 외부 모델
server.Proposal = &proposalsvc.Handler{
    EscrowModel: model.NewEscrowModel(os.Getenv("ESCROW_SERVICE_URL"), &http.Client{}),
}
```

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/externalspec/parser.go` | 신규 — 외부 서비스 스펙 파서 (`@service`, `@endpoint` 파싱) |
| `internal/gluegen/model_impl.go` | DB/외부 모델 분기, 외부 모델 HTTP client 코드젠 |
| `internal/gluegen/gluegen.go` | 모델 수집 시 외부 스펙 로드, 분류 로직 |
| `internal/crosscheck/ddl_coverage.go` | DDL 또는 외부 스펙 존재 확인으로 완화 |
| `internal/crosscheck/ssac_ddl.go` | 외부 모델 타입 DDL 체크 스킵 |
| `internal/crosscheck/func.go` | Func TODO WARNING → ERROR, 금지 import 검사 추가 |
| `internal/funcspec/parser.go` | `FuncSpec.Imports` 필드 추가, import 수집 |
| `internal/orchestrator/validate.go` | 외부 스펙 디렉토리 스캔, 11번째 SSOT 추가 |
| `artifacts/manual-for-ai.md` | 외부 서비스 모델 + @call 순수성 규칙 문법 추가 |

## 의존성

- Phase 020 완료 (crosscheck 정밀도 개선)
- 외부 서비스 스펙 문법 확정 필요

## 미결 사항

1. **외부 스펙 형식** — Go 어노테이션 vs OpenAPI vs 별도 YAML? 위 제안은 Go 어노테이션 (funcspec과 유사한 패턴)
2. **인증** — 외부 서비스 호출 시 API key, OAuth 등 인증 정보 어디서 주입?
3. **에러 처리** — 외부 서비스 타임아웃, 재시도 정책은 어디서 정의?
4. **SSOT 카운트** — 현재 10개 SSOT, 외부 스펙 추가 시 11개로 늘어남. 기존 SSOT에 통합 가능한가?
5. **dummy spec 마이그레이션** — billing @call → @model 전환, transactions 테이블 SSaC 참조 추가

## 검증 방법

```bash
go test ./internal/externalspec/...
go test ./internal/funcspec/...
go test ./internal/crosscheck/...
fullend validate specs/dummy-gigbridge  # 외부 스펙 파싱 + Func 순수성 + 교차 검증
fullend gen specs/dummy-gigbridge artifacts/dummy-gigbridge  # 외부 모델 코드젠
cd artifacts/dummy-gigbridge/backend && go build ./cmd/  # 컴파일 확인
```
