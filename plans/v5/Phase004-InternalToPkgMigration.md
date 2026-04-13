# 🟡 Phase004 — internal 코드젠을 pkg/generate 로 이식 + 즉각 구조 정리 (초벌 완료, 후속 대기)

## 목표

internal 코드젠을 **pkg/generate 에 초벌 이식하면서 동시에 구조적 냄새를 해소** 한다.
이 Phase 의 성공 기준은 **"구조 건전성"** 이지 "internal 산출물 재현" 이 아니다.

- internal 은 **참고용**이며 답습 금지 — 스파게티는 그대로 옮기지 않는다.
- 당장 코드젠 기능이 **부분적으로 불완전해도 OK** — 구조가 건전해야 다음 Phase 의 프로덕션화가 가능하다.
- **복잡성에 따라 if-else / Toulmin 하이브리드** 적용 — 단순은 if-else, 복잡은 Toulmin.
- `internal/` 은 건드리지 않음 (복사 방식).

---

## 전제

- **Phase001 완료** — `pkg/fullend/` 분리, `pkg/parser/ssac.ServiceFunc.Feature` 리네임.
- **Phase002 완료** — `pkg/ground/` 가 generator 요구 반영 신 필드 보유.
- **Phase003 완료** — validate/crosscheck 가 신 필드 기반으로 안정화.

---

## 원칙

### 1. "복사" 는 시작점이지 도달점이 아님

`cp -r` 은 초벌 drop. 그 자리에서 **즉시 구조 정리** 병행. internal 의 구조 그대로 옮기면 실패.

### 2. 복잡성 축에 따른 if-else / Toulmin 판정

**if-else 유지 (단순):**
- 분기 2~3개 이하
- 조건 간 독립 (축 1개)
- 결정이 문장 1~2줄로 표현 가능

**Toulmin 도입 (복잡):**
- 분기 4개 이상
- 조건 교차 (축 2개 이상)
- defeat (우선순위) 의미 필요
- 규칙 집합이 증가형

판정은 이식 중 그 자리에서. 애매하면 Toulmin 선택 (나중에 단순화 쉬움).

internal 의 명백한 Toulmin 대상 (이 Phase 에서 즉시 전환 권장):
- `generate_method_from_iface.go` 의 7-case switch (축 5개 × 7 case)
- `generate_main.go` 의 초기화 블록 조합 (auth/queue/authz/session/cache/file)
- Hurl `build_scenario_order.go` + `classify_mid_step.go` 의 5-phase 우선순위

이 지점은 **복사 + if-else 답습 금지** — Toulmin 으로 즉시 재작성.

### 3. 구조 건전성 > 기능 완성도 (이 Phase 한정)

`go build ./pkg/...` 통과는 필수.
생성 산출물의 `go build` 나 `hurl` 통과는 **목표이지만 blocker 아님**.

### 4. Ground 를 직접 소비

`SymbolTable` 이식 없음. `st.Models[...]`, `st.DDLTables[...]` 접근은 `g.Models[...]`, `g.Tables[...]` 로 기계적 치환. Phase002 에서 확정된 매핑 표 사용.

---

## 즉시 해소 대상 (구조 냄새)

gogin README 의 냄새 5개 중 **이 Phase 에서 처리**:

| 냄새 | 해소 방법 | 난이도 |
|------|----------|--------|
| Flat/Domain 중복 (`*WithDomains` 쌍) | Flat 제거 | 낮음 |
| 매개변수 비대 (`generateMain` 8 params 등) | 입력 struct 로 묶기 | 중간 |
| 결정 분산 (queue init 판정이 7 파일 흩어짐) | 한 곳에 취합, 호출자는 결과만 받음 | 중간 |
| 거대 템플릿 문자열 리터럴 (`main_template.go` 등) | Go `text/template` 분리 또는 구조적 조립 | 중간 |
| 복잡 분기 switch (7-case) | Toulmin 전환 | 높음 |

**즉시 해소 의무**: Flat 중복, 매개변수 비대, 결정 분산 (3개).
**즉시 해소 권장**: 거대 템플릿, 복잡 switch → Toulmin.

---

## 현황

| 위치 | 역할 | 파일 수 | 이식 대상 |
|------|------|--------|----------|
| `internal/gen/gogin/` | Go+Gin 백엔드 생성 | 132 | 복사·정리 → `pkg/generate/gogin/` |
| `internal/gen/react/` | React 프론트엔드 글루 | 23 | 복사·정리 → `pkg/generate/react/` |
| `internal/gen/hurl/` | Hurl smoke 테스트 | 48 | 복사·Toulmin 전환 → `pkg/generate/hurl/` |
| `internal/gen/` 최상위 | `generate.go`, `select_backend.go` | 2 | 복사·정리 → `pkg/generate/` 최상위 |
| `internal/ssac/generator/` | SSaC → Go 핸들러·모델·Handler | 237 | 복사·Ground 치환·정리 → `pkg/generate/gogin/ssac/` |
| `internal/stml/generator/` | STML → React TSX 페이지 | 79 | 복사·정리 → `pkg/generate/react/stml/` |
| `internal/contract/` | 디렉티브·해시·splice·preserve | ~50 | 복사 → `pkg/contract/` |
| `internal/genapi/` | 공유 타입 | 4 | **이식 안 함** — `pkg/fullend.Fullstack` 이 대체 |
| `internal/genmodel/` | 외부 OpenAPI → HTTP 클라이언트 | 40+ | **보류** |

---

## 최종 배치

```
pkg/
├── contract/                      ← internal/contract 복사
├── generate/
│   ├── generate.go                ← internal/gen/generate.go
│   ├── select_backend.go          ← internal/gen/select_backend.go
│   ├── gogin/                     ← internal/gen/gogin (+ Flat 제거, 구조 정리, Toulmin 포인트)
│   │   └── ssac/                  ← internal/ssac/generator (+ Ground 치환)
│   ├── react/                     ← internal/gen/react
│   │   └── stml/                  ← internal/stml/generator
│   └── hurl/                      ← internal/gen/hurl (+ Toulmin ordering)
```

모든 generator 는 `*fullend.Fullstack` + `*rule.Ground` 를 직접 받는다.

### 패키지명 (확정)

- `pkg/generate/gogin/ssac/` — `package ssac`
- `pkg/generate/react/stml/` — `package stml`
- 나머지는 디렉토리명 일치.

### 구조 의도

`gogin`, `react` 는 기술 스택 이름. 수요 기반 확장 옵션 유지 — nisabit.com 운영 데이터에서 Node.js 수요 등 관측 시 `pkg/generate/<stack>/` 로 분기. `select_backend.go` 는 현재 고정 반환, 확장 시점에 `fullend.yaml` 필드와 함께 도입 (YAGNI).

---

## 작업 순서

복사 방식이라 각 Step 비파괴적. 실패 시 pkg 쪽만 삭제하고 재시도.

### Step 1. `internal/contract` → `pkg/contract` 복사

generator 의 디렉티브 해시·preserve 유틸. 선행.

1. `cp -r internal/contract/ pkg/contract/`
2. 패키지 선언 `package contract` 유지.
3. internal 참조 대부분 자립 (hash·splice 계열). import 대체 최소.
4. 구조 정리 기회 식별 — 냄새 있으면 같이 해소.
5. `go build ./pkg/contract/...` + `go test` 통과.
6. 커밋.

### Step 2. `internal/ssac/generator` → `pkg/generate/gogin/ssac` + Ground 치환 + 구조 정리

1. `mkdir -p pkg/generate/gogin`
2. `cp -r internal/ssac/generator/ pkg/generate/gogin/ssac/`
3. 패키지 선언: `package generator` → `package ssac`.
4. Parser import 교체: `internal/ssac/parser` → `pkg/parser/ssac`.
5. **SymbolTable 접근 → Ground 조회 치환** (Phase002 매핑 표 준수):
   - `st.Models[name]` → `g.Models[name]` 등
6. 함수 시그니처: `st *validator.SymbolTable` → `g *rule.Ground`.
7. **구조 정리 병행**:
   - 매개변수 많은 함수 → 입력 struct 로 묶기
   - 중복 제거
   - 판단 분산 → 한 곳에 취합
8. **복잡성 판정 후 Toulmin 대상 전환** (해당 함수만):
   - 분기 4+ 또는 축 2+ 인 함수 식별 → Toulmin 으로 재작성
9. **테스트 16개 (SymbolTable 사용) 을 Ground 기반으로 재작성**. 실패해도 핵심 기능 테스트 우선.
10. `go build ./pkg/...` 통과 확인.

### Step 3. `internal/stml/generator` → `pkg/generate/react/stml` + 구조 정리

1. `mkdir -p pkg/generate/react`
2. `cp -r internal/stml/generator/ pkg/generate/react/stml/`
3. 패키지 선언: `package generator` → `package stml`.
4. Parser import: `internal/stml/parser` → `pkg/parser/stml`.
5. `internal/genapi` 참조는 `pkg/fullend` 로 교체.
6. 냄새 식별 후 즉시 해소 (STML 은 상대적으로 깔끔 — 큰 정리 없을 것).
7. 빌드 통과.

### Step 4. `internal/gen/gogin` → `pkg/generate/gogin` + Flat 제거 + Ground 치환 + 구조 정리 + Toulmin 포인트

본 Phase 최대 단계. 복사·Flat 제거·구조 정리·Toulmin 판정까지 **하나의 일관된 재구성** 으로 취급.

#### Step 4a. 파일 복사 (서브디렉토리 보존)
```
cp -r internal/gen/gogin/. pkg/generate/gogin/
```
주의: Step 2 의 `ssac/` 서브디렉토리 덮어쓰지 않도록.

#### Step 4b. 패키지 선언
`package gogin` 유지.

#### Step 4c. Flat mode 제거

- Flat 전용 파일 삭제 (`generate_server_struct.go`, `generate_main.go` Flat 버전 등).
- `*WithDomains` suffix 제거 → 원래 이름으로 rename.
- `generate.go` 의 Flat/Domain 분기 제거 → 항상 Feature 경로.

#### Step 4d. Ground 치환

Phase002 매핑 표 적용. 매개변수 교체·시그니처 변경·내부 함수도 동일.

#### Step 4e. 구조 정리 (즉시 해소)

**매개변수 비대 해소**:
- `generateMain(artifactsDir, modulePath, queueBackend, sessionBackend, cacheBackend, models, serviceFuncs, policies, fileConfig)` → `generateMain(input MainGenInput)` 형태
- 나머지 8+ 매개변수 함수 동일 처리

**결정 분산 해소**:
- Queue init 판정이 `generate_main.go:68-97` + `collect_subscribers.go` + `has_publish_sequence.go` + `build_queue_blocks.go` 로 분산됨
- → `internal/queue_decision.go` 같은 한 파일로 판정 로직 집약, 호출자는 결과(예: `QueueNeed { Import string; InitBlock string; SubscribeBlock string }`)만 받음

**`main_template.go` 처리**:
- 거대 문자열 리터럴 → `text/template` 분리 또는 구조적 조립
- 복잡성 따라 Toulmin 대상일 수도

#### Step 4f. Toulmin 도입 지점

명백한 Toulmin 필요 지점 식별 후 이 Phase 에서 재작성:

1. **`generate_method_from_iface.go` 의 7-case switch**
   - 축 5개(이름 접두, QueryOpts, 반환 제네릭, 슬라이스, seqType) × 7 case
   - Toulmin warrant 여러 개 + defeat edge 로 재구성

2. **`main.go` 초기화 블록 조합**
   - auth/queue/authz/session/cache/file 6축 독립 조건
   - 단 "필요한 블록만 포함" 이라 각 warrant 단독 판정 가능
   - 실제론 if-else 6개로도 가능하지만, 추가 확장 여지 고려하면 Toulmin 이 장기적으로 유리

판정 기준 경계가 애매하면 **Toulmin 선택** (확장성 우위).

#### Step 4g. Import 경로 교체

- `internal/ssac/parser` → `pkg/parser/ssac`
- `internal/ssac/validator` → **제거** (SymbolTable 사용 없음)
- `internal/ssac/generator` → `pkg/generate/gogin/ssac`
- `internal/genapi` → `pkg/fullend` 또는 제거
- `internal/projectconfig` → `pkg/parser/manifest`
- `internal/contract` → `pkg/contract`
- `internal/policy`, `internal/funcspec`, `internal/statemachine` — internal 유지

### Step 5. `internal/gen/react` → `pkg/generate/react` + 구조 정리

1. `cp -r internal/gen/react/. pkg/generate/react/`
2. `package react` 유지.
3. Import 교체: `internal/genapi` → `pkg/fullend`, `internal/stml/generator` → `pkg/generate/react/stml`, `internal/contract` → `pkg/contract`.
4. 냄새 있으면 해소 (react 도 상대적으로 작음).
5. 빌드 통과.

### Step 6. `internal/gen/hurl` → `pkg/generate/hurl` + Toulmin scenario ordering

1. `cp -r internal/gen/hurl/ pkg/generate/hurl/`
2. `package hurl` 유지.
3. Import 교체.
4. **핵심: `build_scenario_order.go` + `classify_mid_step.go` 의 5-phase + mid-step 순서 로직을 Toulmin 으로 재작성**
   - 5 phase + FK topological sort + state BFS + branch skip 이 한 함수에 압축되어 있음
   - 축 여러 개 + 우선순위 의미 → Toulmin 의 전형적 적용 대상
   - warrant: `IsAuthStep`, `IsPrereqCreate`, `IsStateTransition`, `IsTopLevelCreate`, `IsDelete` 등
   - defeat: `IsDelete.Attacks(IsCreate)` 같은 우선순위 엣지
5. 나머지 step 조립 로직은 if-else 유지 (단순).
6. 빌드 통과.

### Step 7. 최상위 오케스트레이터 복사

1. `cp internal/gen/generate.go pkg/generate/generate.go`
2. `cp internal/gen/select_backend.go pkg/generate/select_backend.go`
3. `package gen` → `package generate`.
4. Import 교체.
5. `internal/gen/` 미삭제.

### Step 8. `orchestrator` 배선 교체

1. `orchestrator/gen_glue.go` — `gen.Generate` 호출을 `pkg/generate.Generate` 로.
2. `orchestrator/gen_ssac.go`, `gen_stml.go` 참조 교체.
3. `orchestrator/parsed.go` — `internal/genapi.ParsedSSOTs` 반환을 `pkg/fullend.Fullstack` + `pkg/rule.Ground` 튜플로 변경.

### Step 9. 최종 빌드 검증

1. `go build ./pkg/... ./internal/... ./cmd/...` 통과 (필수).
2. `go vet` 통과.
3. `go test ./pkg/ground/...` 통과 (필수 — 기반이라).
4. `go test ./pkg/...` 부분 실패 허용.
5. orchestrator 가 pkg/generate 만 호출하는지 grep 확인.

---

## 주의사항

### R1. Ground 신 필드 미충분 발견 시

Phase002 설계로 generator 의 모든 접근을 못 소화하면:
- 소규모 누락: Phase002 로 복귀해 필드 확장 후 재개.
- 대규모 누락: 본 Phase 보류 + Phase002 재설계.

### R2. 웹 도메인 의미 보존

`pkg/parser/manifest/deploy.go:7` 의 `Deploy.Domain string` 는 웹 도메인. 리네임 대상 아님.

### R3. 타입 분리

`pkg/parser/ssac.ServiceFunc` (Feature 필드) ≠ `internal/ssac/parser.ServiceFunc` (Domain 필드). orchestrator 가 pkg 로 배선된 뒤엔 internal 쪽 타입 무호출이라 문제 없음.

### R4. Step 4 원자성 권장

복사·Flat 제거·Ground 치환·구조 정리·Toulmin 포인트가 한 Step 에 묶이면 중간 상태 빌드 깨짐. **논리적 원자성 유지 위해 단일 커밋** 권장. 실패 시 `git reset --hard`.

다만 구조 정리가 길어지면 **Step 4 내부를 여러 sub-커밋** 으로 쪼개도 OK:
- 4a~4d (복사·Flat 제거·Ground 치환) 커밋 — 기본 구조 이식
- 4e (구조 정리) 커밋 — 매개변수·분산 정리
- 4f (Toulmin 포인트) 커밋 — 복잡 분기 재작성
각 커밋마다 빌드 통과 보증이라면 분할 OK.

### R5. 테스트 재작성 범위

ssac generator 의 SymbolTable 사용 테스트 16개. Ground 기반으로 재작성. 전부 통과가 의무는 아님 — 본 Phase 의 "기능 불완전 OK" 원칙 적용. 단 핵심 로직 테스트는 통과해야.

### R6. funcspec / policy / statemachine / projectconfig 미이식

pkg/generate 쪽 복사본은 internal 경로 계속 참조. 차후 별도 Phase.

### R7. ffignore 미수정

프로젝트 규약. filefunc 경보 유발 시 코드 조정으로 해소.

### R8. "불완전 허용" 의 경계

| 항목 | 허용 |
|------|------|
| `go build ./pkg/...` 실패 | **불가** |
| `go vet` 경고 | 허용 |
| `go test ./pkg/ground/...` 실패 | **불가** |
| `go test ./pkg/generate/...` 실패 | 허용 (이식 중간) |
| 생성 `artifacts/backend/go build` 실패 | 허용 |
| 생성 `hurl --test` 실패 | 허용 |
| 핸들러 body 비어있음 | 허용 |
| 산출물 파일 일부 누락 | 허용 |

---

## 검증 방법

### 정적 (필수)
1. `go build ./pkg/... ./internal/... ./cmd/...` 통과.
2. `go vet` 통과.
3. `go test ./pkg/ground/...` 통과.

### 배선 확인
- `grep -rn "internal/gen\"\|internal/genapi\|internal/ssac/generator\|internal/stml/generator" internal/orchestrator/` 결과 0.

### 어휘 확인
- `grep -rn "\.Domain\b" pkg/generate/` 결과에 웹 도메인 외 잔존 없음.

### 구조 건전성 (Phase005 에서 본격 측정)

본 Phase 에선 **체감 검토**만.
- 매개변수 많은 함수 남아있는가?
- 결정 분산 해소됐는가?
- Toulmin 적용 지점이 의도된 곳인가?

---

## 완료 조건 (Definition of Done)

- [ ] `pkg/contract/` 존재
- [ ] `pkg/generate/{gogin, gogin/ssac, react, react/stml, hurl}` + 최상위 2파일 존재
- [ ] pkg/generate/gogin 에서 Flat mode 관련 파일 삭제 (suffix `_with_domains` 제거)
- [ ] generator 의 SymbolTable 매개변수가 Ground 로 치환됨
- [ ] orchestrator 가 pkg/generate 만 호출
- [ ] orchestrator/parsed.go 가 pkg/fullend.ParseAll + pkg/ground.Build 기반
- [ ] **매개변수 비대 해소** (기존 8+ params 함수 → struct 로 묶음)
- [ ] **결정 분산 해소** (queue init 등 한 곳 취합)
- [ ] **Toulmin 적용** — method_from_iface switch, main 초기화 블록, hurl 시나리오 순서 중 최소 2곳
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] `go test ./pkg/ground/...` 통과
- [ ] `internal/` 은 변경되지 않은 상태로 유지됨
- [ ] 커밋 메시지: `pkg/generate 이식 + 구조 정리 + Toulmin 도입 — internal 원본 유지`

---

## 다음 Phase 예고

- **Phase005** — 실용 검증 + 구조 건전성 지표 측정.
- **Phase00N** — internal 일괄 삭제 (pkg 안정화 후).
- **Phase00M+** — 프로덕션화 (기능 완성, 성능, 에러 처리 등).
