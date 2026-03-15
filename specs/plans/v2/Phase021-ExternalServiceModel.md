# ✅ Phase 021: 패키지 접두사 모델 판단 규칙 + Func 순수성 강제

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
| External | fullend 내장 (SSOT 불필요) | `*http.Client` 기반 구현체 | Phase 022 대상 |
| Session | fullend 내장 (SSOT 불필요) | `SessionModel` 인터페이스 + 백엔드별 구현체 | Phase 022 대상 |
| Cache | fullend 내장 (SSOT 불필요) | `CacheModel` 인터페이스 + 백엔드별 구현체 | Phase 022 대상 |
| File | fullend 내장 (SSOT 불필요) | `LocalFileModel` + `S3Model` 인터페이스 + 구현체 | Phase 022 대상 |
| Queue | 미정 | 미정 | 향후 |

Phase 021에서는 **패키지 접두사 기반 모델 판단 규칙**과 **Func 순수성 강제**를 구현.
내장 모델(Session, Cache, File, External)은 Phase 022에서 구현.

### B. Func 순수성 강제 + TODO 레벨 수정

Phase 020에서 이관된 항목:

| # | 카테고리 | 문제 |
|---|---|---|
| 1 | Func TODO 레벨 | 본체 미구현이 WARNING이지만 ERROR가 맞음 |
| 2 | Func 순수성 검증 | @call func이 I/O 패키지를 import하면 ERROR |

외부 모델이 제공되어야 Func 순수성 강제가 의미 있음 (I/O를 @model로 대체할 수 있으니까).

## 설계

### 1. 모델 판단 규칙 (패키지 접두사 기반)

SSaC `@model` 문법에서 **패키지 접두사 유무**로 즉시 분기. 폴백 체인 없음.

| SSaC 문법 | 판단 | validate |
|---|---|---|
| `@model Gig.FindByID(...)` | 접두사 없음 → DDL 테이블 | DDL + interface 검증 |
| `@model session.Session.Get(...)` | `session` 패키지 → import 경로 탐색 | interface 검증 |
| `@model cache.Cache.Get(...)` | `cache` 패키지 → import 경로 탐색 | interface 검증 |
| `@model file.File.Upload(...)` | `file` 패키지 → import 경로 탐색 | interface 검증 |
| `@model escrow.Escrow.Hold(...)` | `escrow` 패키지 → import 경로 탐색 | interface 검증 |
| `@model payment.Payment.Charge(...)` | `payment` 패키지 → 사용자 커스텀 | interface 검증 |
| `@model mydb.MyDB.Query(...)` | `mydb` 패키지 → 사용자 커스텀 DB | interface 검증 (DDL 검증 없음 — 사용자 감수) |

**규칙:**
- **접두사 없음** → DB 모델 (DDL이 SSOT)
- **접두사 있음** → 해당 패키지 경로에서 Go interface 파싱 → 교차 검증
  - interface 있음 → 메서드/파라미터/리턴 타입 검증
  - interface 없음 → WARNING ("검증 불가")
  - 패키지 자체 없음 → ERROR

fullend/pkg/session, pkg/cache, pkg/file이든 사용자 커스텀이든 **동일한 검증 경로**. 특별 취급 없음.

### 2. 코드젠 분기

```
모델 수집 → 접두사 유무로 분류
  ├── 접두사 없음 → DB 모델 → 기존 model_impl.go 로직 (SQL 기반, fullend 코드젠)
  └── 접두사 있음 → 패키지 모델 → interface 파싱 (코드젠은 패키지가 제공)
```

**DB 모델 (접두사 없음 — fullend 코드젠):**
```go
// SSaC: @model Gig.FindByID(...)
type GigModel interface {
    FindByID(id int64) (*Gig, error)
}
type gigModelImpl struct { db *sql.DB }
```

**패키지 모델 (접두사 있음 — 패키지가 구현 제공):**
```go
// SSaC: @model escrow.Escrow.Hold(...)
// escrow 패키지의 interface → ssac validate가 파싱하여 교차 검증 (수정지시서 016)
// 구현체는 패키지 자체가 제공 (fullend 코드젠 대상 아님)
```

### 3. Func 순수성 강제

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

### 4. Func TODO 레벨 수정

본체 미구현(`// TODO: implement`만 있는 stub)이 `WARNING`으로 분류됨. 실행 시 빈 값을 반환하므로 `ERROR`가 맞음.

**수정**: func.go:83 `Level: "WARNING"` → `Level: "ERROR"`

### 5. 교차 검증 변경

| 규칙 | 현재 | 변경 |
|---|---|---|
| @model → DDL | 모든 @model에 DDL 테이블 요구 | 접두사 없는 @model만 DDL 요구, 접두사 있으면 DDL 스킵 |
| @result ↔ DDL | 접두사 있는 타입도 DDL 체크 | 접두사 있는 @model의 @result는 DDL 체크 스킵 |
| DDL → SSaC | DDL 테이블 미참조 WARN | 변경 없음 |
| Func 순수성 | 없음 (신규) | 금지 import 감지 시 ERROR |
| Func TODO | WARNING | ERROR |

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/ddl_coverage.go` | 접두사 있는 @model은 DDL 체크 스킵 |
| `internal/crosscheck/ssac_ddl.go` | 접두사 있는 @model의 @result DDL 체크 스킵 |
| `internal/crosscheck/func.go` | Func TODO WARNING → ERROR, 금지 import 검사 추가 |
| `internal/funcspec/parser.go` | `FuncSpec.Imports` 필드 추가, import 수집 |
| `artifacts/manual-for-ai.md` | 패키지 접두사 모델 규칙 + @call 순수성 규칙 문법 추가 |

## 의존성

- Phase 020 완료 (crosscheck 정밀도 개선)
- SSaC 수정지시서 016 (패키지 접두사 @model 파싱 + Go interface 교차 검증)

## 검증 방법

```bash
go test ./internal/funcspec/...
go test ./internal/crosscheck/...
fullend validate specs/dummy-gigbridge  # Func 순수성 + 교차 검증
```
