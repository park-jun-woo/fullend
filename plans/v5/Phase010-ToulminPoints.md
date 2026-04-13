# ✅ Phase010 — 결정 지점 구조 정비 — Decide\* 순수 함수 수렴 (완료: aab3c48)

## 기준점 (2026-04-13 추가)

**"if-else 2-depth 이내에 해결 안되면 Toulmin"** — 사용자 확정 기준.

- `depth 1`: `if X {} else if Y {} else {}` (flat chain / switch case)
- `depth 2`: 위의 각 분기 내부에 중첩 if 1회
- `depth 3+`: Toulmin 적용 대상

**조건식 AND/OR 복합도는 depth 에 불포함** (수평 확장). depth 는 수직 중첩만 계수.

## 3 포인트 depth 재평가

| 포인트 | 최소 표현 depth | 판정 |
|-------|--------------|------|
| B1 — method_from_iface dispatch | switch(1) + default 내 if(2) = **2** | ✗ **Toulmin 제외** — `DecideMethodPattern` 순수 함수 |
| B2 — main 초기화 블록 | 축별 독립 if = **1** | ✗ **Toulmin 제외** — `DecideMainInit` 순수 함수 |
| B3 — hurl mid classifier | early-return chain = **2** | ✗ **Toulmin 제외** — `DecideMidStepClass` 순수 함수 |

> 결론: **3 포인트 모두 2-depth 이내**로 해결 가능. Phase010 은 "Toulmin 도입"이 아니라 **판정 로직을 `Decide*` 순수 함수로 수렴 + facts/needs/decision struct 분리** 의 구조 정비 Phase 로 축소한다.

> Toulmin graph 는 depth 초과 또는 규칙 집합이 증가형·defeat 의미 필요한 경우에만. 현 3 포인트는 그런 속성이 없다.

## 목표 (수정)

Phase009 구조 정리 위에서 **복잡한 결정 지점 3군데를 `Decide*` 순수 함수로 수렴**.

1. **method_from_iface** 의 7-case switch → `DecideMethodPattern(facts) Pattern`
2. **main 초기화 블록 조합** → `DecideMainInit(facts) InitNeeds`
3. **hurl mid classifier** → `DecideMidStepClass(step, maps) MidDecision`

각 `Decide*` 는:
- 입력: 평범 facts struct
- 출력: enum/struct 단일 결정
- 내부: `switch` 또는 `if-else` 2-depth 이내
- 호출자: 반환값을 평범 `switch` 로 소비 (판정 로직 無)

성공 기준:
- 3군데 모두 `Decide*` 로 분리
- 판정 로직이 호출자에 남지 않음
- 행동 보존 (bit-동일 아니어도 기능 동일)
- `go build` + `go test ./pkg/...` 통과

---

## 전제

- **Phase009 완료** — 구조 정리 선행.
- internal/gen 은 dead code 유지.

---

## 설계 선행 (Phase002 같은 체크포인트)

Toulmin 도입 3군데 각각 **설계 문서** 선행:

- `plans/v5/Phase010-MethodDispatchDesign.md`
- `plans/v5/Phase010-MainInitDesign.md`
- `plans/v5/Phase010-ScenarioOrderDesign.md`

각 설계 문서:
1. 현재 internal 로직 정밀 분석 (분기 조건 전수)
2. 축 목록 (orthogonal dimensions)
3. warrant 목록 + defeat edges
4. 반환 struct 정의
5. 호출자 소비 패턴

**설계 검토 후 구현 착수**. 설계 오류 발견 시 재설계.

---

## 원칙 (재확인)

- **2-depth 이내는 if-else, 초과는 Toulmin** (상단 기준점).
- "결정 함수 하나 → Pattern/Needs/Decision 반환" 구조 유지.
- 호출자는 평범한 Go switch 로 소비.
- Toulmin 적용 불필요 시에도 **facts/decision struct 분리**와 **판정 로직 수렴**은 그대로 수행.

---

## 범위

### Part B1. method_from_iface.go 의 7-case switch

- 위치: `pkg/generate/gogin/generate_method_from_iface.go`
- 축 5개 × 7 case → **단일 `switch` + default 내 if = 2-depth**
- 구현: `DecideMethodPattern(facts MethodFacts) Pattern` 순수 함수
- 결과: `Pattern` 반환 → 호출자 평범 switch 로 구현 선택
- 상세: `Phase010-MethodDispatchDesign.md`

### Part B2. main 초기화 블록 조합

- 위치: `pkg/generate/gogin/generate_main.go`
- 6축 독립 → **축별 독립 if, depth 1**
- 구현: `DecideMainInit(facts MainFacts) InitNeeds` 순수 함수
- 결과: `InitNeeds{Auth, Queue, Authz, Session, Cache, File, NeedsContextImport}` 반환
- 상세: `Phase010-MainInitDesign.md`

### Part B3. hurl mid classifier

- 위치: `pkg/generate/hurl/classify_mid_step.go`
- 7 분기 early-return chain → **2-depth 이내**
- 구현: `DecideMidStepClass(step, maps) MidDecision` 순수 함수
- 5-phase 배치 / FK topo / state BFS / auth order / prereq split 은 **그대로 유지** (알고리즘 교체 대상 아님)
- 상세: `Phase010-ScenarioOrderDesign.md`

---

## 작업 순서

### Step 1. 설계 문서 3개 작성

Part B1/B2/B3 각각 설계 문서. 사용자 리뷰 후 구현 착수.

### Step 2. Part B1 — method_from_iface

- 설계 문서 → `decide_method_pattern.go` 신설
- 순수 switch (depth 2 이내) → `Pattern` 반환
- `generate_method_from_iface.go` 의 switch 를 Pattern 소비로 교체
- 기존 테스트 유지

### Step 3. Part B2 — main 초기화

- 설계 문서 → `decide_main_init.go` 신설
- `InitNeeds` 반환 (depth 1 축별 판정 함수)
- `generate_main.go` 가 InitNeeds 소비해 블록 조립

### Step 4. Part B3 — hurl mid classifier

- 설계 문서 → `decide_mid_step_class.go` 신설
- 순수 switch-case (depth 2 이내) → `MidDecision` 반환
- `build_scenario_order.go` 가 `decideMidStepClass` 호출로 교체. 기존 `classify_mid_step.go` 제거

### Step 5. 검증

- `go build` + `go vet` + `go test ./pkg/...`
- `fullend gen dummys/gigbridge/specs /tmp/x` — 실행 에러 없음
- `Decide*` 단위 테스트 추가 (분기 전수 표)

---

## 주의사항

### R1. 설계 선행 필수

3군데 모두 설계 문서 리뷰 후 구현. "구현하면서 설계" 금지.

### R2. Toulmin 판정 기준 엄격

기준점: **if-else 2-depth 초과 시 Toulmin**. 본 Phase 3 포인트는 모두 2-depth 이내로 해결되어 Toulmin 미적용. "내 눈엔 복잡해 보임" 같은 이유로 Toulmin 추가 금지.

보조 기준 (depth 이내여도 Toulmin 고려 가능한 예외):
- 규칙 집합이 명시적 증가형 (evidence 기반 defeat 필요)
- 판정 과정 trace 가 감사/디버깅에 필수

### R3. 행동 변화 허용 범위

- 생성된 산출물의 **실질 기능 동일** 유지
- 산출물의 공백·순서 등 사소한 차이는 허용

### R4. 부분 완료 허용

시간 초과 시 Part B1 (가장 명확) 만이라도 완료하고 B2, B3 는 별도 Phase 로 분리 가능.

---

## 완료 조건 (Definition of Done)

- [x] Part B1/B2/B3 설계 문서 작성 및 기준점(2-depth) 재평가 반영
- [ ] method_from_iface switch → `DecideMethodPattern` 순수 함수 전환
- [ ] main 초기화 블록 조합 → `DecideMainInit` 순수 함수 전환
- [ ] hurl mid classifier → `DecideMidStepClass` 순수 함수 전환
- [ ] `go build` + `go vet` + `go test ./pkg/...` 통과
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 성공
- [ ] 커밋 메시지: `refactor(generate): Decide* 순수 함수로 결정 로직 수렴 (3포인트)`

---

## 다음 Phase

- **Phase011** — Dummy 실용 검증 + 구조 건전성 지표 측정 (internal 대비 개선 확인).
- **Phase00N (별도)** — internal/gen, internal/genapi, internal/ssac/generator, internal/stml/generator, internal/contract 일괄 삭제.
