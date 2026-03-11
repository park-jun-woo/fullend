# Phase 021: 외부 서비스 모델 지원 + Func 순수성 강제

## 배경

현재 `@model`은 DDL 테이블(DB)에 1:1 종속. 외부 서비스 통신은 `@call func`으로만 가능.
`@call func`을 순수 로직(I/O 금지)으로 제한하기로 결정했으므로, 외부 서비스 I/O를 담당할 계층이 필요.

### @model vs @call 역할 구분 (확정)

- **@model** = 모든 I/O 추상화. 서비스 코드에서 I/O가 필요하면 반드시 @model을 통해 접근.
- **@call func** = 순수 비즈니스 로직. I/O 금지. 계산, 판단, 변환만 허용.
- 서비스 시퀀스가 @call(계산) → @model(I/O) 순서를 조율.

### 모델 I/O 분류 (6종)

| # | 분류 | 목적 | 예시 |
|---|---|---|---|
| 1 | **DB** | 영속 데이터 | PostgreSQL, MySQL, MongoDB, DynamoDB |
| 2 | **Session** | 사용자 귀속 상태 (로그인, 장바구니 등) | Redis, DB |
| 3 | **Cache** | 빈번한 요청 효율화 (TTL 기반 임시 저장) | Redis, Memcached, CDN |
| 4 | **File** | 파일/오브젝트 스토리지 | S3, GCS, 로컬 파일시스템 |
| 5 | **External** | 외부 서비스 API | Stripe, Twilio, SendGrid |
| 6 | **Queue** | 비동기 메시지 | Kafka, RabbitMQ, SQS |

Session과 Cache는 기술적으로 같은 인프라(Redis 등)를 쓸 수 있지만, 목적이 다르므로 모델 계층에서는 분리.

## 현재 모델 흐름

```
DDL SQL → 테이블 정의
  ↓
SSaC @model Gig.FindByID(...) → 인터페이스 메서드
  ↓
fullend gluegen → *sql.DB 기반 구현체 자동 생성
```

## 목표

### A. 모델 I/O 종류별 확장

| 종류 | 정의 SSOT | 코드젠 결과 | 현재 상태 |
|---|---|---|---|
| DB | DDL SQL | `*sql.DB` 기반 구현체 | ✅ 구현 완료 |
| External | 외부 서비스 OpenAPI | `*http.Client` 기반 구현체 | Phase 021 대상 |
| Session | 미정 | 미정 | 향후 |
| Cache | 미정 | 미정 | 향후 |
| File | 미정 | 미정 | 향후 |
| Queue | 미정 | 미정 | 향후 |

Phase 021에서는 **External** 모델을 우선 구현. 나머지는 필요 시 후속 Phase에서 추가.

### B. Func 순수성 강제 + TODO 레벨 수정

Phase 020에서 이관된 항목:

| # | 카테고리 | 문제 |
|---|---|---|
| 1 | Func TODO 레벨 | 본체 미구현이 WARNING이지만 ERROR가 맞음 |
| 2 | Func 순수성 검증 | @call func이 I/O 패키지를 import하면 ERROR |

외부 모델이 제공되어야 Func 순수성 강제가 의미 있음 (I/O를 @model로 대체할 수 있으니까).

## 설계

### 1. 외부 서비스 스펙 SSOT

외부 서비스가 공개하는 **OpenAPI 문서를 그대로 SSOT로 사용**. 새 문법을 발명하지 않음.
대부분의 SaaS (Stripe, Twilio, SendGrid 등)가 OpenAPI 스펙을 공개 제공하며, fullend는 이미 OpenAPI 파서를 보유.

```
specs/<project>/
├── api/openapi.yaml              ← 우리 API (기존)
├── external/
│   ├── escrow.openapi.yaml       ← 외부 결제 서비스 OpenAPI
│   └── notification.openapi.yaml ← 외부 알림 서비스 OpenAPI (예시)
```

외부 OpenAPI에서 fullend가 추출하는 정보:
- **operationId** → 모델 메서드명 매핑 (e.g. `createEscrowHold` → `Escrow.CreateHold`)
- **request schema** → 메서드 파라미터 타입
- **response schema** → 메서드 리턴 타입
- **서버 URL** → base URL (환경변수로 오버라이드 가능)
- **security scheme** → API key, Bearer token 등 인증 방식

장점:
- 새 문법 발명 불필요 — 기존 OpenAPI 파서 재활용
- 외부 서비스 업데이트 시 OpenAPI 파일만 교체하면 모델 자동 갱신
- 교차 검증 시 SSaC @model 파라미터를 외부 OpenAPI request schema와 대조 가능

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

**ERROR 메시지 지침:**
```
[ERROR] Func ↔ SSaC: <funcName> — @call func에서 I/O 패키지 %q import 금지. @call func은 순수 계산/판단 로직만 허용됩니다. DB, 네트워크, 파일 등 I/O가 필요하면 @model을 활용하세요.
```

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
| `internal/gluegen/external_model.go` | 신규 — 외부 OpenAPI에서 HTTP client 모델 코드젠 |
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

1. **외부 스펙 형식** — ✅ 결정: 외부 서비스가 공개하는 OpenAPI 문서를 그대로 사용. 새 문법 발명 불필요. `specs/<project>/external/*.openapi.yaml`에 배치
2. **인증** — 외부 서비스 호출 시 API key, OAuth 등 인증 정보 어디서 주입?
3. **에러 처리** — 외부 서비스 타임아웃, 재시도 정책은 어디서 정의?
4. **SSOT 카운트** — 외부 OpenAPI는 기존 OpenAPI 파서로 처리하므로 별도 SSOT 종류 추가 없이 external/ 디렉토리 스캔만 추가
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
