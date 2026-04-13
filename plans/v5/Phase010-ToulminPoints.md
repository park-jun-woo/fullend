#  Phase010 — Toulmin 포인트 도입 (3군데)

## 목표

Phase009 구조 정리 위에서 **복잡한 결정 지점 3군데를 Toulmin graph 로 전환**.

1. **method_from_iface** 의 7-case switch (축 5개)
2. **main 초기화 블록 조합** (6축 독립 조건)
3. **hurl 시나리오 순서** (5-phase + FK topological + state BFS + branch skip)

Toulmin 도입은 "**결정을 명시적 규칙 그래프로 표현** + **defeat 관계로 우선순위 선언**" 이 목적. 단순 if/for 로 충분한 곳엔 도입 금지.

성공 기준:
- 3군데 모두 Toulmin graph 적용
- 각 지점별 "Decide* 함수 → Pattern/Info/Order struct 반환" 구조 (호출자는 평범 Go 소비)
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

- **단순은 if-else, 복잡은 Toulmin** (축 3+ 직교 또는 우선순위 의미 있을 때 Toulmin).
- Toulmin 도입 시에도 "결정 함수 하나 → 패턴 반환" 구조 유지.
- 호출자는 평범한 Go if/for 로 Pattern 소비.

---

## 범위

### Part B1. method_from_iface.go 의 7-case switch

- 위치: `pkg/generate/gogin/generate_method_from_iface.go`
- 축 5개 (이름 접두, QueryOpts 유무, 반환 제네릭, 슬라이스 여부, seqType) × 7 case

**Warrant 초안**:
- `IsWithTx` (defeater, 최상위 우선)
- `IsCursorPaginated`, `IsOffsetPaginated`, `IsSliceReturn`
- `IsFind`, `IsCreate`, `IsUpdate`, `IsDelete`

**Defeat edges**: WithTx > Pagination > Slice > Find/Create/...

**결과**: `Pattern` 타입 반환 → 호출자가 평범한 switch 로 구현 선택.

### Part B2. main 초기화 블록 조합

- 위치: `pkg/generate/gogin/generate_main.go`
- 6축 독립 조건 (auth/queue/authz/session/cache/file)

**Warrant**:
- `NeedsAuth`, `NeedsQueue`, `NeedsAuthz`, `NeedsSession`, `NeedsCache`, `NeedsFile`

**Defeat edges**: 없음 (축이 독립적).

**결과**: `InitNeeds { Auth, Queue, Authz, Session, Cache, File bool + ... 상세 }` 반환.

렌더링은 평범한 Go if/for (struct 필드 소비).

### Part B3. hurl 시나리오 순서

- 위치: `pkg/generate/hurl/build_scenario_order.go` + `classify_mid_step.go`
- 5-phase (Auth / Prereq / Mid / Read / Delete)
- FK topological sort + state BFS + branch skip

**Warrant**:
- `IsAuthStep`
- `IsTopLevelCreate`, `IsStateTransition`, `IsNestedCreate`
- `IsUpdate`, `IsRead`, `IsDelete`
- `HasFKDependency` (defeat으로 순서 조정)

**Toulmin Loop 후보**: 복잡 루프 조건 (state BFS).
**Subscribe 후보**: 단계별 이벤트 분산 (phase 전이).

**결과**: 정렬된 `[]Step`.

---

## 작업 순서

### Step 1. 설계 문서 3개 작성

Part B1/B2/B3 각각 설계 문서. 사용자 리뷰 후 구현 착수.

### Step 2. Part B1 — method_from_iface

- 설계 문서 → `decide_method_pattern.go` 신설
- Toulmin graph + Evaluate → Pattern 반환
- `generate_method_from_iface.go` 의 switch 를 Pattern 소비로 교체
- 기존 테스트 유지

### Step 3. Part B2 — main 초기화

- 설계 문서 → `decide_main_init.go` 신설
- warrant 6개 + InitNeeds 반환
- `generate_main.go` 가 InitNeeds 소비해 블록 조립

### Step 4. Part B3 — hurl 시나리오 순서

- 설계 문서 → `decide_scenario_order.go` 신설
- Loop 또는 계층 graph
- Subscribe 적용 여부 평가
- `build_scenario_order.go` + `classify_mid_step.go` 의 로직을 새 구조로 이전

### Step 5. 검증

- `go build` + `go vet` + `go test ./pkg/...`
- `fullend gen` — 실행 에러 없음
- Toulmin warrant 단위 테스트 추가 (각 warrant 별 true/false 판단)

---

## 주의사항

### R1. 설계 선행 필수

3군데 모두 설계 문서 리뷰 후 구현. "구현하면서 설계" 금지.

### R2. Toulmin 판정 기준 엄격

이 3군데 외 "내 눈엔 복잡해 보임" 같은 이유로 추가 Toulmin 적용 금지. 기준:
- 분기 4+ 직교 축
- 우선순위/defeat 의미 필요
- 규칙 집합이 증가형

### R3. 행동 변화 허용 범위

- 생성된 산출물의 **실질 기능 동일** 유지
- 산출물의 공백·순서 등 사소한 차이는 허용

### R4. 부분 완료 허용

시간 초과 시 Part B1 (가장 명확) 만이라도 완료하고 B2, B3 는 별도 Phase 로 분리 가능.

---

## 완료 조건 (Definition of Done)

- [ ] Part B1/B2/B3 설계 문서 작성 및 리뷰 통과
- [ ] method_from_iface switch → Toulmin graph 전환
- [ ] main 초기화 블록 조합 → Toulmin
- [ ] hurl 시나리오 순서 → Toulmin
- [ ] `go build` + `go vet` + `go test ./pkg/...` 통과
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 성공
- [ ] 커밋 메시지: `refactor(generate): Toulmin 포인트 3군데 도입`

---

## 다음 Phase

- **Phase011** — Dummy 실용 검증 + 구조 건전성 지표 측정 (internal 대비 개선 확인).
- **Phase00N (별도)** — internal/gen, internal/genapi, internal/ssac/generator, internal/stml/generator, internal/contract 일괄 삭제.
