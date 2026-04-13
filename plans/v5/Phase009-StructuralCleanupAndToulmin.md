# Phase009 — 구조 정리 + Toulmin 포인트 도입

## 목표

Phase008 까지 pkg/generate 가 동작하는 상태가 됐다. 이제 **구조 건전성을 실질적으로 올리는** 작업.

1. **매개변수 비대 해소** — `generateMain` 등 다매개변수 함수를 struct 로.
2. **결정 분산 해소** — Queue init 등 판정 로직을 한 곳에 취합.
3. **거대 템플릿 처리** — `main_template.go` 같은 큰 문자열 리터럴을 `text/template` 로 분리.
4. **Toulmin 포인트 3군데 도입**:
   - `generate_method_from_iface.go` 의 7-case switch → Toulmin warrant + defeat
   - `main` 초기화 블록 조합 (auth/queue/authz/session/cache/file) → Toulmin
   - hurl 시나리오 순서 결정 (5-phase + FK topological) → Toulmin Loop + Subscribe

성공 기준:
- 지표가 internal 대비 개선 (Phase010 측정)
- 행동 보존 — 생성 산출물 기능 유지 (단 bit-동일 의무 없음; 구조적 개선 허용)
- `go build` + `go test ./pkg/...` 통과

---

## 전제

- **Phase008 완료** — pkg/generate 가 실제 동작, orchestrator 가 pkg 호출.
- internal/gen 은 dead code 로 남음 (본 Phase 에서도 유지).

---

## 원칙 (재확인)

- **단순은 if-else, 복잡은 Toulmin** (축 3+ 직교 또는 우선순위 의미 있을 때 Toulmin).
- 구조 정리는 순수 리팩토링 (행동 변화 최소화).
- Toulmin 도입 시에도 "결정 함수 하나 → 패턴 반환" 구조 유지 (호출자는 평범한 Go 로 소비).

---

## 범위

### 포함 (Part A — 구조 정리, 순수 리팩토링)

1. **매개변수 비대 해소**
   - `generateMain` (8 params) → `generateMain(input MainGenInput)`
   - 기타 7+ params 함수 전수 조사 후 변환
   
2. **결정 분산 해소**
   - Queue init 판정 (`generate_main.go` + `collect_subscribers.go` + `has_publish_sequence.go` + `build_queue_blocks.go` 분산) → `queue_decision.go` 한 곳에 취합, 결과 struct (`QueueNeed{Import, InitBlock, SubscribeBlock}`) 반환
   - 기타 유사 패턴 식별 후 처리

3. **템플릿 분리**
   - `main_template.go`, `main_with_domains_template.go`, `query_opts_template.go` 같은 거대 리터럴
   - Go `text/template` 기반 `.tmpl` 파일 + 렌더 함수로 분리
   - 복잡성 높은 것은 Toulmin 대상 (아래)

### 포함 (Part B — Toulmin 포인트 3군데)

**B1. method_from_iface.go 의 7-case switch**
- 축 5개 (이름 접두, QueryOpts 유무, 반환 제네릭, 슬라이스 여부, seqType) × 7 case
- Toulmin warrant 들:
  - `IsWithTx` (defeater, 최상위 우선)
  - `IsCursorPaginated`, `IsOffsetPaginated`, `IsSliceReturn`
  - `IsFind`, `IsCreate`, `IsUpdate`, `IsDelete`
- Defeat edges: 우선순위 표현 (WithTx > Pagination > Slice > Find/Create/...)
- 결과: `Pattern` 타입 반환 → 호출자가 switch 로 구현 선택

**B2. main 초기화 블록 조합**
- 6축 독립 조건 (auth/queue/authz/session/cache/file)
- 각 축별 warrant 단독 판정 (defeat 없음)
- 결과: `InitNeeds { Auth, Queue, Authz, Session, Cache, File bool + ... 상세 }` 반환
- 렌더링은 평범한 Go if/for (warrant 결과 소비)

**B3. hurl 시나리오 순서**
- 5-phase (Auth / Prereq / Mid / Read / Delete)
- FK topological sort + state BFS + branch skip
- Toulmin Loop 후보 (복잡 루프 조건) + Subscribe 후보 (단계별 이벤트)
- 결과: 정렬된 `[]Step` 반환

### 포함하지 않음

- 나머지 Toulmin 도입 — 본 Phase 3군데만. 추가는 별도 Phase.
- Dummy 실용 검증 — Phase010.

---

## 작업 순서

### Step 1. 구조 지표 현 상태 캡처

Phase010 에서 측정할 지표의 "before" 값 수집 (내부 함수 매개변수 평균, 결정 분산 파일 수 등). 본 Phase 가 얼마나 개선했는지 추후 비교용.

### Step 2. Part A — 구조 정리 (순수 리팩토링)

각 작업 독립 커밋:
- `refactor(generate): generateMain 을 MainGenInput struct 로 수렴`
- `refactor(generate): queue init 판정 queue_decision.go 로 집약`
- `refactor(generate): main_template 을 text/template 분리`

기타 작업은 발견 시 유사 패턴으로 독립 커밋.

### Step 3. Part B1 — method_from_iface Toulmin 전환

- Toulmin graph 설계 (warrant 목록, defeat edges)
- `decide_method_pattern.go` 신설 — Graph + Evaluate → Pattern 반환
- `generate_method_from_iface.go` 의 switch 를 Pattern 소비로 교체
- 테스트: 기존 동작 유지

### Step 4. Part B2 — main 초기화 Toulmin

- `decide_main_init.go` 신설 — warrant 6개 + InitNeeds 반환
- `generate_main.go` 가 InitNeeds 소비해 블록 조립

### Step 5. Part B3 — hurl 시나리오 순서 Toulmin

- `decide_scenario_order.go` 신설 — Loop 또는 계층 graph
- Subscribe 적용 여부 평가 (Step 단위 이벤트 분산)
- `build_scenario_order.go` + `classify_mid_step.go` 의 로직을 새 구조로 이전

### Step 6. 검증

- `go build` + `go vet` + `go test ./pkg/...`
- `fullend gen dummys/gigbridge/specs /tmp/x` — 실행 에러 없음
- 생성 산출물 diff — internal 대비 개선되었거나 동등 (회귀 아님)

---

## 주의사항

### R1. 매 커밋 단위로 빌드 통과

구조 정리와 Toulmin 도입 모두 점진적. 한 커밋에 변경 최소화.

### R2. Toulmin 판정 기준 엄격 적용

Part B 3군데 외 "내 눈엔 복잡해 보임" 같은 이유로 추가 Toulmin 적용 금지. 기준:
- 분기 4+ 직교 축
- 우선순위/defeat 의미 필요
- 규칙 집합이 증가형

모호하면 if-else 유지. Phase010 지표 확인 후 Phase00X 에서 추가 고려.

### R3. 행동 변화 허용 범위

- 내부 구조 자유롭게 개선
- 생성된 산출물의 **실질 기능 동일** 유지 (API 시그니처, 비즈니스 로직)
- 산출물의 공백·순서 등 사소한 차이는 허용

### R4. Phase009 완료 시점 판단

- 3개 Toulmin 포인트 모두 적용 + 구조 정리 Part A 완료 → DoD
- 시간 초과 시 Part B1 (가장 명확) 만이라도 완료하고 Part B2, B3 는 별도 Phase 로 분리 가능

---

## 완료 조건 (Definition of Done)

- [ ] 매개변수 8+ 함수가 pkg/generate 에 존재하지 않음
- [ ] Queue init 판정이 단일 파일에 집약
- [ ] 거대 템플릿 리터럴 1개 이상 text/template 분리
- [ ] method_from_iface 의 switch → Toulmin graph 전환
- [ ] main 초기화 블록 조합 → Toulmin
- [ ] hurl 시나리오 순서 → Toulmin
- [ ] `go build` + `go vet` + `go test ./pkg/...` 통과
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 성공
- [ ] 커밋: `refactor(generate): 구조 정리 + Toulmin 포인트 3군데 도입`

---

## 다음 Phase

- **Phase008** — Dummy 실용 검증 + 구조 건전성 지표 측정 (internal 대비 개선 확인).
- **Phase00N (별도)** — internal/gen, internal/genapi, internal/ssac/generator, internal/stml/generator, internal/contract 일괄 삭제.
