# Phase001 — internal 코드젠을 pkg/generate 로 이식

## 목표

`internal/` 에 흩어져 있던 **코드젠 로직 전체를 `pkg/generate/` 로 이식**하여 `go build ./...` 가 통과하도록 한다.
이식은 **순수 이동 + 어휘 통일 + 레거시 제거** 까지만. Toulmin 기반 리팩토링은 이 단계에서 **하지 않는다**.

## 전략 요약

- 행동 변화 금지: 이식 전후 **생성되는 산출물은 bit-동일**해야 한다 (`diff == 0`).
- `Domain` → `Feature` 일괄 리네임 (단, 웹 도메인 의미는 보존).
- Flat mode 제거 (파서 단에서 금지되는 코드 경로라 사용 불가능).
- `internal/genmodel` 은 본 Phase에서 **이식 보류** (외부 API 클라이언트 독립 기능).

검증 기준:
- `go build ./pkg/... ./internal/... ./cmd/...` 통과
- `go vet` 통과
- dummy 프로젝트(gigbridge) 기준 `fullend gen` 실행 시 이전 산출물과 diff 0

---

## 현황

| 위치 | 역할 | 파일 수 | 이식 대상 |
|------|------|--------|----------|
| `internal/gen/gogin/` | Go+Gin 백엔드 생성 | 132 | O |
| `internal/gen/react/` | React 프론트엔드 글루 | 23 | O |
| `internal/gen/hurl/` | Hurl smoke 테스트 | 48 | O |
| `internal/gen/` | 최상위 오케스트레이터 (`generate.go`, `select_backend.go`) | 2 | O |
| `internal/ssac/generator/` | SSaC → Go 핸들러·모델 인터페이스·Feature Handler 구조체 | 237 | O |
| `internal/stml/generator/` | STML → React TSX 페이지 | 79 | O |
| `internal/genapi/` | 공유 타입 (`ParsedSSOTs`, `GenConfig`, `STMLGenOutput`, `Backend`) | 4 | O |
| `internal/genmodel/` | 외부 OpenAPI → HTTP 클라이언트 | 40+ | **보류** (Phase002 또는 이후) |
| `pkg/generate/` | (삭제 완료) | 0 | — |

---

## 최종 배치

```
pkg/generate/
├── generate.go                    ← internal/gen/generate.go
├── select_backend.go              ← internal/gen/select_backend.go
├── api/                           ← internal/genapi/
│   ├── parsed_ssots.go
│   ├── gen_config.go
│   ├── stml_gen_output.go
│   └── backend.go
├── backend/                       ← internal/gen/gogin/
├── frontend/                      ← internal/gen/react/
├── hurl/                          ← internal/gen/hurl/
├── ssac/                          ← internal/ssac/generator/
└── stml/                          ← internal/stml/generator/
```

`api` 서브패키지로 두는 이유: 최상위에 두면 `pkg/generate` 가 자기 타입 `api.ParsedSSOTs` 를 import 해야 하는 순환 위험이 있음. 독립 서브패키지가 안전.

---

## 변경 파일 목록

### 삭제 (이식 후)

- `internal/gen/gogin/**`
- `internal/gen/react/**`
- `internal/gen/hurl/**`
- `internal/gen/generate.go`, `internal/gen/select_backend.go`
- `internal/ssac/generator/**`
- `internal/stml/generator/**`
- `internal/genapi/**`

### 추가 (이식된 파일들)

`pkg/generate/api/` — 4개 파일
`pkg/generate/backend/` — 132개 → **Flat mode 제거 후 ~95개 예상**
`pkg/generate/frontend/` — 23개
`pkg/generate/hurl/` — 48개
`pkg/generate/ssac/` — 237개
`pkg/generate/stml/` — 79개
`pkg/generate/generate.go`, `pkg/generate/select_backend.go`

### 수정 (호출자 측 import 교체)

- `internal/orchestrator/gen_glue.go` — `gen.Generate` 호출 → `pkg/generate.Generate`
- `internal/orchestrator/gen_ssac.go` — `ssacgenerator` import 경로 교체
- `internal/orchestrator/gen_stml.go` — `stmlgenerator` import 경로 교체
- `internal/orchestrator/parsed.go` — `genapi.ParsedSSOTs` → `pkg/generate/api.ParsedSSOTs`
- `internal/orchestrator/*.go` 중 `genapi` 참조하는 모든 파일

---

## 작업 순서

### Step 1. `internal/genapi` → `pkg/generate/api` 이식

가장 독립성 높은 타입 패키지부터. 순환 의존성 리스크 없음.

1. 파일 복사: `cp -r internal/genapi/ pkg/generate/api/`
2. 패키지 선언 변경: `package genapi` → `package api`
3. 호출자 import 교체:
   - `"github.com/park-jun-woo/fullend/internal/genapi"` → `"github.com/park-jun-woo/fullend/pkg/generate/api"`
   - alias: 필요 시 `genapi "…/pkg/generate/api"` 로 유지 가능 (호출자 변경 최소화)
4. `internal/genapi/` 삭제
5. `go build ./pkg/... ./internal/...` 통과 확인

### Step 2. `internal/ssac/generator` → `pkg/generate/ssac` 이식

SSaC generator 는 다른 generator 가 의존하는 모델 인터페이스를 만드는 선행 단계이므로 먼저.

1. 디렉토리 이동: `mv internal/ssac/generator pkg/generate/ssac`
2. 패키지 선언 변경: `package generator` → `package ssac` (또는 `package ssacgen` — 충돌 시)
3. Parser import 교체:
   - `"github.com/park-jun-woo/fullend/internal/ssac/parser"` → `"github.com/park-jun-woo/fullend/pkg/parser/ssac"`
   - `"github.com/park-jun-woo/fullend/internal/ssac/validator"` → (실제 필요 여부 확인 후) 해당 pkg 경로로
4. `funcspec` 참조 — 현재 `internal/funcspec` 그대로 유지할지 `pkg/parser/funcspec` 으로 전환할지 확인. **Phase 1에서는 internal 유지** (부담 최소화).
5. 호출자(`orchestrator/gen_ssac.go`) import 교체.
6. `go build` 통과 확인.

### Step 3. `internal/stml/generator` → `pkg/generate/stml` 이식

SSaC와 독립이지만 orchestrator 가 둘 다 호출하므로 같은 패턴.

1. 디렉토리 이동: `mv internal/stml/generator pkg/generate/stml`
2. 패키지 선언: `package generator` → `package stml` (또는 `stmlgen`).
3. Parser import 교체: `internal/stml/parser` → `pkg/parser/stml`.
4. 호출자(`orchestrator/gen_stml.go`) import 교체.
5. 빌드 확인.

### Step 4. `internal/gen/gogin` → `pkg/generate/backend` 이식 + Flat 제거 + Feature 리네임

가장 방대한 단계. 세 개 작업을 **한 커밋 안에** 묶는다. 이유: 부분 이식 시 중간 상태가 빌드 안 됨.

#### Step 4a. 파일 이동
```
mv internal/gen/gogin pkg/generate/backend
```

#### Step 4b. 패키지 선언 변경
`package gogin` → `package backend`

#### Step 4c. Flat mode 제거

삭제 대상 파일 (파서가 금지하는 경로):
- `generate_server_struct.go` (`…_with_domains` 가 대체)
- `generate_main.go` (`…_with_domains` 가 대체)
- `transform_service_files.go` (`…_with_domains` 가 대체)
- `auth.go` (Flat 전용 `CurrentUser`)
- `generate_auth_stub.go` (있다면)
- `main_template.go` (Flat 템플릿)
- 기타 `has_domains.go` 로 분기되는 Flat 브랜치 내부 코드

`generate.go` 의 `if hasDomains(…) { ... } else { ... }` 분기 → else 제거, **항상 Feature 경로** 로 직진.

#### Step 4d. Domain → Feature 리네임

타깃:
- 구조체 필드: `ServiceFunc.Domain` 은 `pkg/parser/ssac` 에 있으므로 **별도 확인** (아래 주의 섹션 참조). 필요 시 `Feature` 로 리네임 + 전체 참조 동기화.
- 함수명:
  - `hasDomains` → `hasFeatures`
  - `uniqueDomains` → `uniqueFeatures`
  - `generateDomainHandler` → `generateFeatureHandler`
  - `collectModelsForDomain` → `collectModelsForFeature`
  - `collectFuncsForDomain` → `collectFuncsForFeature`
  - `domainNeedsDB` → `featureNeedsDB`
  - `domainNeedsAuth` → `featureNeedsAuth`
  - `domainNeedsJWTSecret` → `featureNeedsJWTSecret`
  - `anyDomainNeedsAuth` → `anyFeatureNeedsAuth`
  - `generateServerStructWithDomains` → `generateServerStruct` (Flat 제거로 suffix 불필요)
  - `generateMainWithDomains` → `generateMain`
  - `transformServiceFilesWithDomains` → `transformServiceFiles`
  - `generateAuthStubWithDomains` → `generateAuthStub`
  - `generate_main_with_domains.go` → `generate_main.go`
  - 기타 `*_with_domains.go` suffix 제거
- 변수명: `domain` → `feature`, `domains` → `features`, `domainDir` → `featureDir`, `domainHandler` → `featureHandler` 등
- 주석: 모두 교체

#### Step 4e. Import 경로 교체
- `internal/ssac/parser` → `pkg/parser/ssac`
- `internal/ssac/validator` → 해당 pkg 경로
- `internal/genapi` → `pkg/generate/api`
- `internal/policy` → (유지) — pkg 이식은 별도 Phase
- `internal/projectconfig` → `pkg/parser/manifest` (확인 필요)
- `internal/funcspec` → (유지)
- `internal/statemachine` → (유지)

#### Step 4f. 호출자 배선 교체
- `internal/gen/select_backend.go` 의 `return &gogin.GoGin{}` → `return &backend.Backend{}` (타입명도 같이 갱신).

### Step 5. `internal/gen/react` → `pkg/generate/frontend` 이식

1. `mv internal/gen/react pkg/generate/frontend`
2. `package react` → `package frontend`
3. Import 경로 교체:
   - `internal/genapi` → `pkg/generate/api`
4. 빌드 확인.

### Step 6. `internal/gen/hurl` → `pkg/generate/hurl` 이식

1. `mv internal/gen/hurl pkg/generate/hurl`
2. `package hurl` → `package hurl` (동일)
3. Import 교체:
   - `internal/genapi` → `pkg/generate/api`
   - `internal/statemachine`, `internal/policy` → (유지)
4. 빌드 확인.

### Step 7. 최상위 오케스트레이터 이식

1. `internal/gen/generate.go` + `internal/gen/select_backend.go` → `pkg/generate/generate.go`, `pkg/generate/select_backend.go`
2. 내부 import:
   - `internal/gen/gogin` → `pkg/generate/backend`
   - `internal/gen/react` → `pkg/generate/frontend`
   - `internal/gen/hurl` → `pkg/generate/hurl`
   - `internal/genapi` → `pkg/generate/api`
3. `internal/gen/` 디렉토리 삭제 (`internal/gen/generate.go`, `select_backend.go` 외 남은 파일 없어야 함).

### Step 8. `orchestrator` 배선 교체

1. `orchestrator/gen_glue.go`:
   - `"github.com/park-jun-woo/fullend/internal/gen"` → `"github.com/park-jun-woo/fullend/pkg/generate"`
   - `gen.Generate(…)` → `generate.Generate(…)`
2. `orchestrator/gen_ssac.go`:
   - `ssacgenerator "github.com/park-jun-woo/fullend/internal/ssac/generator"` → `ssacgen "github.com/park-jun-woo/fullend/pkg/generate/ssac"`
3. `orchestrator/gen_stml.go`:
   - `stmlgenerator "github.com/park-jun-woo/fullend/internal/stml/generator"` → `stmlgen "github.com/park-jun-woo/fullend/pkg/generate/stml"`
4. `orchestrator/parsed.go` 및 `genapi` 참조 전부 → `pkg/generate/api`.

### Step 9. 최종 빌드·검증

1. `go build ./pkg/... ./internal/... ./cmd/...` 통과
2. `go vet ./pkg/... ./internal/... ./cmd/...` 통과
3. 누락 import·dead 참조 제거
4. dummy 프로젝트(`dummys/gigbridge`) 에서 `fullend gen` 실행
5. 기존 artifact snapshot 과 diff 0 확인

---

## 주의사항

### ServiceFunc.Feature 리네임 범위

`pkg/parser/ssac/service_func.go` 의 `Domain string` 필드가 `Feature` 로 바뀌면:
- `pkg/parser/ssac/parse_dir_entry.go` — 할당부 동시 변경
- `pkg/validate/ssac/*` — 이미 `Domain` 참조가 있을 수 있음 (확인 필요)
- `pkg/crosscheck/*` — 동일
- `pkg/ground/*` — 동일
- `test_parse_domain_folder_test.go` → `test_parse_feature_folder_test.go` (파일명도)

**파서 쪽 필드 리네임은 본 Phase 에서 같이 수행** (코드젠만 바꾸면 필드명 불일치로 빌드 깨짐).

### 웹 도메인 의미 보존

`pkg/parser/manifest/deploy.go:7` 의 `Deploy.Domain string` 는 **웹 도메인 주소** 를 뜻하므로 **리네임 대상 아님**.

리네임 작업은 반드시 다음 방식 중 하나로:
- 파일 단위 `Edit` 로 개별 확인 후 변경
- `go refactor` / IDE 리네임 기능
- `sed` 는 위험 — 전역 치환은 웹 도메인 의미까지 오염

### 패키지명 충돌

- `internal/gen/hurl` 의 `package hurl` → `pkg/generate/hurl` 의 `package hurl` 로 유지 가능 (충돌 없음).
- `internal/ssac/generator` 의 `package generator` 는 `pkg/generate/ssac` 로 옮기면서 `package ssacgen` 로 변경 권장 (패키지명 일반성 낮추기).
- `internal/stml/generator` → `package stmlgen` 동일 이유.
- `internal/gen/gogin` 의 `package gogin` → `pkg/generate/backend` 는 `package backend` 로.

### `internal/gen/` 디렉토리

이식 완료 후 `internal/gen/` 디렉토리는 **완전 삭제**. `rmdir internal/gen` 성공해야 함 (남은 파일 없음 의미).

### funcspec / policy / statemachine / projectconfig 위치

본 Phase 에서는 **이식하지 않는다**. pkg 이식은 별도 Phase (필요 시). 본 Phase 의 generator 들은 `internal/*` 의 이 패키지들을 계속 참조.

### 테스트 파일

`test_*_test.go` 전부 동반 이식. 패키지명·import 경로 업데이트 필요. 테스트 실패가 회귀 감지 핵심 수단.

### pkg/ssac, pkg/stml 충돌 확인

`pkg/` 루트에 `ssac`, `stml` 디렉토리가 이미 있다면 충돌. `pkg/generate/ssac`, `pkg/generate/stml` 는 중첩 경로라 충돌 없음. 확인 필요.

---

## 의존성

- 없음 (내부 리팩토링).
- `pkg/parser/ssac`, `pkg/parser/stml`, `pkg/parser/manifest` 가 이미 존재하고 사용 가능한 상태여야 함 (이미 존재 확인됨).

---

## 검증 방법

### 정적 검증
1. `go build ./pkg/... ./internal/... ./cmd/...` 통과 (scripts/ 는 기존 이슈로 제외).
2. `go vet ./pkg/... ./internal/... ./cmd/...` 통과.
3. `grep -rn "internal/gen/" internal/ pkg/` 결과 0건 (모든 참조 이주 완료).
4. `grep -rn "internal/genapi" internal/ pkg/` 결과 0건.
5. `grep -rn "internal/ssac/generator" internal/ pkg/` 결과 0건.
6. `grep -rn "internal/stml/generator" internal/ pkg/` 결과 0건.

### 런타임 검증
1. `dummys/gigbridge/` 에서 `fullend gen` 실행.
2. 이식 전 artifact snapshot 과 `diff -r artifacts-before/ artifacts-after/` 결과 0 바이트 차이.
3. `artifacts/backend/` 에서 `go build ./...` 통과.
4. `artifacts/tests/smoke.hurl` 이 문법적으로 유효 (`hurl --test --dry-run`).

### 어휘 검증
1. `grep -rn "hasDomains\|Domain\b" pkg/generate/` 에서 웹 도메인 관련 외 참조 0건.
2. `pkg/parser/ssac/service_func.go` 의 필드가 `Feature` 로 리네임됐는지 확인.
3. `specs/service/<feature>/*.ssac` 의 주석·에러 메시지에 "feature" 용어 통일.

---

## 완료 조건 (Definition of Done)

- [ ] `internal/gen/` 디렉토리 삭제됨
- [ ] `internal/genapi/` 디렉토리 삭제됨
- [ ] `internal/ssac/generator/` 디렉토리 삭제됨
- [ ] `internal/stml/generator/` 디렉토리 삭제됨
- [ ] `pkg/generate/{api,backend,frontend,hurl,ssac,stml}` + 최상위 파일 존재
- [ ] Flat mode 관련 파일 삭제됨 (`*_with_domains` 접미 제거 포함)
- [ ] `Domain` → `Feature` 리네임 완료 (웹 도메인 제외)
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] gigbridge dummy artifact diff 0
- [ ] 커밋 메시지: `pkg/generate 이식 완료 — internal/gen·ssac·stml 이동, Flat 제거, Domain→Feature`

---

## 다음 Phase 예고

- **Phase002** — dummy 회귀 스위트 구축 (gigbridge + zenflow 자동 diff 검증 스크립트).
- **Phase003** — 복잡 로직 3군데를 Toulmin 기반으로 리팩토링:
  1. `backend/generate_method_from_iface.go` 의 7-case switch (축 5개)
  2. `backend/generate_main.go` 의 초기화 블록 조합 (auth/queue/authz/session/cache/file)
  3. `hurl/build_scenario_order.go` + `classify_mid_step.go` 의 5-phase 우선순위 그래프
