# Phase 022: fullend 내장 모델 (External, Session, Cache, File)

## 배경

Phase 021에서 패키지 접두사 기반 모델 판단 규칙이 확립됨.
접두사가 있으면 해당 패키지 경로에서 Go interface를 파싱하여 교차 검증.
fullend/pkg든 사용자 커스텀이든 동일한 규칙.

Phase 022에서는 fullend가 기본 제공하는 내장 모델 패키지를 구현.
사용은 선택 — 사용자가 자기 패키지로 대체 가능.

## 대상

| 종류 | 패키지 | 구현체 |
|---|---|---|
| External | `pkg/external` | OpenAPI 기반 HTTP client |
| Session | `pkg/session` | Redis, DB, Memory |
| Cache | `pkg/cache` | Redis, Memcached, Memory |
| File | `pkg/file` | LocalFile, S3 |

## 설계

### 1. External 모델 (pkg/external)

외부 서비스가 공개하는 OpenAPI 문서를 SSOT로 사용. fullend는 이미 OpenAPI 파서를 보유.

**OpenAPI 배치:**
```
specs/<project>/
├── api/openapi.yaml              ← 우리 API (기존)
├── external/
│   ├── escrow.openapi.yaml       ← 외부 결제 서비스 OpenAPI
│   └── notification.openapi.yaml ← 외부 알림 서비스 OpenAPI
```

**OpenAPI에서 추출하는 정보:**
- **operationId** → 모델 메서드명 매핑 (e.g. `createEscrowHold` → `Escrow.CreateHold`)
- **request schema** → 메서드 파라미터 타입
- **response schema** → 메서드 리턴 타입
- **서버 URL** → base URL (환경변수로 오버라이드 가능)
- **security scheme** → API key, Bearer token 등 인증 방식

**코드젠 결과:**
```go
type EscrowModel interface {
    Hold(req EscrowHoldRequest) (*EscrowHoldResponse, error)
}
type escrowModelImpl struct {
    baseURL string
    client  *http.Client
}
```

**SSaC에서의 사용:**
```go
// @post HoldResponse result = escrow.Escrow.Hold({GigID: gig.ID, Amount: gig.Budget})
```

**공식 외부 서비스 패키지:**
- 메이저 서비스(Stripe, Google, Slack, Twilio, SendGrid, Telegram 등)는 공개 OpenAPI를 활용하여 fullend가 공식 패키지로 제공
- 별도 레지스트리(`fullend-ext/`)로 분리하여 코어와 독립 버전 관리
- `fullend ext-gen <openapi.yaml>` 도구로 사용자도 직접 패키지 생성 가능

**미결:**
- 인증 정보 주입 방법 (API key, OAuth 등)
- 타임아웃, 재시도 정책 정의 위치

### 2. Session 모델 (pkg/session)

key-value + TTL 구조. 사용자 귀속 상태 (로그인, 장바구니 등).

**인터페이스:**
```go
type SessionModel interface {
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Get(ctx context.Context, key string) (string, error)
    Delete(ctx context.Context, key string) error
}
```

**SSaC에서의 사용:**
```go
// @get Session session = session.Session.Get({token: request.Token})
// @post Session result = session.Session.Set({token: userID, TTL: 3600})
// @delete Session result = session.Session.Delete({token: request.Token})
```

**fullend.yaml 설정:**
```yaml
session:
  backend: redis  # redis | db | memory
```

### 3. Cache 모델 (pkg/cache)

Session과 동일한 key-value + TTL 구조. 목적만 다름 (데이터 효율화).

**인터페이스:**
```go
type CacheModel interface {
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Get(ctx context.Context, key string) (string, error)
    Delete(ctx context.Context, key string) error
}
```

**SSaC에서의 사용:**
```go
// @get CachedGig gig = cache.Cache.Get({key: "gig:" + request.ID})
// @post CacheResult result = cache.Cache.Set({key: "gig:" + gig.ID, value: gig, TTL: 300})
// @delete CacheResult result = cache.Cache.Delete({key: "gig:" + request.ID})
```

**fullend.yaml 설정:**
```yaml
cache:
  backend: redis  # redis | memcached | memory
```

### 4. File 모델 (pkg/file)

key(경로) → 파일 바이너리. 2종 구현체 기본 제공.

**인터페이스:**
```go
type FileModel interface {
    Upload(ctx context.Context, key string, body io.Reader) error
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
}
```

**구현체:**
- `LocalFileModel` — 로컬 파일시스템 (`os.Create`, `os.Open`)
- `S3Model` — AWS S3 (`s3.PutObject`, `s3.GetObject`)

**SSaC에서의 사용:**
```go
// @post FileResult result = file.File.Upload({key: path, body: request.File})
// @get FileData data = file.File.Download({key: path})
// @delete FileResult result = file.File.Delete({key: path})
```

**fullend.yaml 설정:**
```yaml
file:
  backend: s3  # s3 | local
  s3:
    bucket: my-bucket
    region: ap-northeast-2
  local:
    root: ./uploads
```

S3가 기본. 비용이 거의 무시할 수준이고, 로컬은 서버 이전/스케일아웃 시 문제됨.

## 공통 원칙

- 모든 내장 모델은 패키지 접두사 규칙을 따름 (Phase 021 규칙과 동일)
- fullend가 특별 취급하지 않음 — 사용자 커스텀 패키지와 동일한 검증 경로
- 사용자가 대체 시 fullend.yaml 설정은 무시됨 (사용자 패키지가 자체 설정 관리)

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `pkg/external/external.go` | 신규 — OpenAPI 기반 ExternalModel 코드젠 + HTTP client 구현체 |
| `pkg/session/session.go` | 신규 — SessionModel 인터페이스 + Redis/DB/Memory 구현체 |
| `pkg/cache/cache.go` | 신규 — CacheModel 인터페이스 + Redis/Memcached/Memory 구현체 |
| `pkg/file/file.go` | 신규 — FileModel 인터페이스 + LocalFile/S3 구현체 |
| `internal/projectconfig/config.go` | session/cache/file backend 설정 파싱 |
| `internal/gluegen/model_impl.go` | 내장 모델 의존성 주입 코드젠 |

## 의존성

- Phase 021 완료 (패키지 접두사 모델 판단 규칙)

## 미결 사항

1. **외부 서비스 공식 패키지** — Stripe, Google, Slack 등 fullend 기본 제공 ext 패키지는 별도 레지스트리(`fullend-ext/`)로 분리 검토
2. **fullend 독립 조직** — `fullend.org` 도메인 확보 가능. `geul-org` → `fullendhq` 이관 시점
3. **Queue 모델** — Kafka, RabbitMQ, SQS 등. 필요 시 후속 Phase

## 검증 방법

```bash
go test ./pkg/session/...
go test ./pkg/cache/...
go test ./pkg/file/...
fullend validate specs/dummy-gigbridge  # 내장 모델 interface 교차 검증
fullend gen specs/dummy-gigbridge artifacts/dummy-gigbridge
cd artifacts/dummy-gigbridge/backend && go build ./cmd/
```
