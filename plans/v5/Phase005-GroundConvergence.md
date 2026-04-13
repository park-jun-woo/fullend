# ✅ Phase005 — 내부 함수 Ground 수렴 (완료)

## 목표

Phase004 에서 **stub 으로 남겨둔** pkg/generate 의 내부 함수들을 **`rule.Ground` 기반으로 수렴** 한다. Phase002 에서 구축한 Ground 신 필드(`Models`, `Tables`, `Ops`, `ReqSchemas`)를 직접 소비하도록 재배선.

성공 기준:
- `go build ./pkg/... ./internal/... ./cmd/...` 통과
- `go test ./pkg/...` 통과 (기존 테스트 + 재작성된 generator 테스트)
- pkg/generate/gogin 내부 함수가 `*rule.Ground` 받음 (SymbolTable 의존 제거)
- ServiceFunc 타입이 `pkg/parser/ssac.ServiceFunc` 로 통일

**본 Phase 에선 pkg/generate.Generate 의 stub 을 해소하지 않는다** — 내부 함수만 정돈. stub → 실제 배선은 Phase006.

---

## 전제

- **Phase004 완료** — pkg/generate 에 파일 이식 완료, build 통과, stub 상태.
- pkg/generate/gogin, gogin/ssac, react, react/stml, hurl 모두 존재.
- 내부 함수들이 아직 `internal/ssac/validator.SymbolTable` 및 `internal/ssac/parser.ServiceFunc` 를 참조 중 (또는 pkg/parser/ssac 쪽과 혼합).

---

## 범위

### 포함

- pkg/generate/gogin/**/*.go 내부 함수 시그니처 전환
- pkg/generate/gogin/ssac/**/*.go 동일 전환
- pkg/generate/hurl/**/*.go 에서 SymbolTable/ServiceFunc 참조 정돈
- 테스트 fixture 재작성 (SymbolTable 구성 → Ground 구성)

### 포함하지 않음

- Flat mode 제거 — Phase006
- pkg/generate.Generate stub 해소 및 orchestrator 배선 교체 — Phase006
- 매개변수 비대·결정 분산 구조 정리 — Phase007
- Toulmin 포인트 도입 — Phase007
- Dummy 실용 검증 — Phase008

---

## 작업 순서

### Step 1. 타입 매핑 정밀 조사

사전 조사 필요:
- 어떤 함수가 `*validator.SymbolTable` 을 받는가? (전수 grep)
- 어떤 함수가 `[]parser.ServiceFunc` (internal) 을 받는가?
- Ground 의 어떤 필드로 1:1 대응되는가?

산출: 매핑 표 (함수명 → 대체 접근). Phase002 Ground 설계 문서와 cross-check.

### Step 2. gogin/ssac 내부 Ground 치환

가장 먼저. SSaC generator 가 SymbolTable 에 가장 의존적.

- `st.Models[name]` → `g.Models[name]`
- `st.DDLTables[name]` → `g.Tables[name]`
- `st.Operations[id]` → `g.Ops[id]`
- `st.RequestSchemas[id]` → `g.ReqSchemas[id]`
- 함수 시그니처: `st *validator.SymbolTable` → `g *rule.Ground`

순환 의존성 리스크: pkg/generate/gogin/ssac 는 pkg/rule 참조 OK (rule 은 순수 타입).

### Step 3. gogin 본체 Ground 치환

gogin 의 `generate_method_from_iface`, `analyze_http_func` 등 ssac generator 가 호출하는 지점과 gogin 자체가 직접 SymbolTable 을 쓰는 지점.

### Step 4. hurl 내부 정돈

hurl 은 SymbolTable 직접 의존 적음. 주로 `[]ServiceFunc`, `OpenAPIDoc`, `StateDiagrams` 소비.
타입 정렬만: `internal/ssac/parser.ServiceFunc` → `pkg/parser/ssac.ServiceFunc` 등.

### Step 5. react 내부 정돈

react 는 STML 및 OpenAPI 만 쓰므로 영향 작음.

### Step 6. 테스트 재작성

SymbolTable 을 직접 구성하던 테스트들을 `rule.Ground` 구성 방식으로 재작성.
ssac generator 의 16개 SymbolTable 기반 테스트가 주 대상.

### Step 7. 빌드 + 기존 테스트 전수 통과

- `go build ./pkg/... ./internal/... ./cmd/...`
- `go vet`
- `go test ./pkg/...`

---

## 주의사항

### R1. pkg 가 internal 을 import 하는 상태 유지

Phase005 에서도 `internal/ssac/validator`, `internal/funcspec`, `internal/policy` 등은 그대로 참조 가능 (Go internal 규칙 허용). Ground 수렴으로 이들 의존을 일부 제거할 수 있으나 완전 제거는 Phase006 이후.

### R2. 테스트 실패 범위

기존 ssac generator 테스트 16개는 SymbolTable 직접 구성 방식이라 Ground 로 재작성 시 일시적으로 빌드 실패 가능. 테스트 재작성 커밋을 별도 분리.

### R3. Ground 신 필드 부족 발견

내부 함수가 요구하는 정보가 Ground 에 없으면 — Phase002 로 되돌아가 필드 추가. 예:
- DDL.index.{table} 정보 필요 시 TableInfo 확장
- 기타 누락 발견 시 설계 업데이트

### R4. 커밋 단위

- Step 2 (ssac Ground 치환): 1 커밋
- Step 3 (gogin 본체): 1 커밋
- Step 4 (hurl): 1 커밋
- Step 5 (react): 1 커밋 (영향 작으면 Step 4 와 통합 가능)
- Step 6 (테스트 재작성): 1 커밋

---

## 완료 조건 (Definition of Done)

- [ ] pkg/generate/gogin/ssac/**/*.go 에서 `validator.SymbolTable` 참조 0건
- [ ] pkg/generate/gogin/**/*.go 에서 `validator.SymbolTable` 참조 0건
- [ ] 내부 함수 시그니처가 `*rule.Ground` 받음
- [ ] ServiceFunc 타입이 `pkg/parser/ssac.ServiceFunc` 로 통일 (internal/ssac/parser 참조 제거)
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go test ./pkg/...` 전수 통과
- [ ] 커밋 메시지: `refactor(generate): 내부 함수 Ground 수렴 — SymbolTable 의존 제거`

---

## 다음 Phase

- **Phase006** — Flat mode 제거 + pkg/generate.Generate stub 해소 + orchestrator 배선 교체.
- **Phase007** — 구조 정리 (매개변수·결정 분산·템플릿) + Toulmin 포인트 3군데.
- **Phase008** — Dummy 실용 검증.
