# Phase010 Part B3 — ScenarioOrderDesign

> `pkg/generate/hurl/classify_mid_step.go` 의 Mid classifier 를 **`MidDecision` struct + `DecideMidStepClass` 순수 함수**로 수렴하는 설계.
> 기준점 (2-depth) 재평가 반영 — Toulmin graph 불채택.
> 5-phase 배치 / FK topo / state BFS / auth order / prereq split 등 정렬 알고리즘은 **그대로 유지**.

---

## 1. 현행 구조 분석

### 1.1 엔트리포인트

`buildScenarioOrder(doc, specsDir, diagrams, serviceFuncs) []scenarioStep`
— `pkg/generate/hurl/build_scenario_order.go:24~114`

호출 체인: `generate.Generate` → `generateHurlTests` → `buildScenarioOrder`.

### 1.2 5-phase 조립

| Phase | 입력 | 정렬 | 결합 순서 | 위치 |
|-------|------|------|---------|------|
| Prereq | Mid 중 `order<0` ∧ auth FK 일치 | 원본 order | 최전방 | `splitPrereqSteps` L99 |
| Auth | `isAuthOperation` | `authOrder()` register=0, login=1 | Prereq 다음 | L79 |
| Mid | POST/PUT/PATCH | `order ASC → depth ASC → path ASC` | Auth 다음 | L84, `classifyMidStep` |
| Read | GET | depth ASC → path ASC | Mid 다음 | `sortByDepthPath` L95 |
| Delete | DELETE | FK topological (children→parents) | 최후미 | `sortDeletesByFK` L97 |

**이 phase 배치 자체는 선형 concat, depth 1.** 변경 대상 아님.

### 1.3 Mid classifier — 현행 `classify_mid_step.go:8~29`

```go
func classifyMidStep(s scenarioStep, stateOps, branchSkip map[string]bool,
    transitionOrder map[string]int, resourceFirstTransition map[string]float64,
) (float64, bool) {
    if stateOps[s.OperationID] {
        if branchSkip[s.OperationID] {        // depth 2
            return 0, false
        }
        return float64(transitionOrder[s.OperationID]), true
    }
    if s.Method != "POST" {
        return 900.0, true
    }
    parent := findParentResource(s.Path)
    if parent == "" {
        return -1.0, true
    }
    if ft, ok := resourceFirstTransition[parent]; ok {
        return ft + 0.5, true
    }
    return -0.5, true
}
```

**최대 depth = 2** (stateOps 안의 branchSkip). **기준점 이내** → **Toulmin 제외**.

### 1.4 현행 문제

- 반환값 `(float64, bool)` — 의미가 숨겨짐 (class 명시 없음)
- 호출자가 `order` 만 받고 **어떤 클래스로 분류됐는지** 로그/테스트가 어려움
- early-return chain 자체는 가독 양호

---

## 2. 설계

### 2.1 반환 타입

```go
// pkg/generate/hurl/decide_mid_step_class.go
type StepClass int

const (
    ClassExcluded             StepClass = iota   // state branch skip
    ClassStateTransition                         // order = transitionOrder[id]
    ClassUpdate                                  // order = 900.0
    ClassTopLevelCreate                          // order = -1.0
    ClassNestedUnderTransition                   // order = firstTransition + 0.5
    ClassNestedOrphan                            // order = -0.5
)

type MidDecision struct {
    Class   StepClass
    Order   float64
    Include bool            // ClassExcluded → false, 그 외 true
}

type StepFacts struct {
    Step                    scenarioStep
    IsStateOp               bool        // stateOps[id]
    IsBranchSkip            bool        // branchSkip[id]
    TransitionOrder         int         // transitionOrder[id]
    ParentResource          string      // findParentResource(path)
    FirstTransition         float64     // resourceFirstTransition[parent]
    HasFirstTransition      bool
}
```

### 2.2 판정 함수

```go
func DecideMidStepClass(f StepFacts) MidDecision {
    if f.IsStateOp {
        if f.IsBranchSkip {
            return MidDecision{Class: ClassExcluded, Include: false}
        }
        return MidDecision{
            Class: ClassStateTransition,
            Order: float64(f.TransitionOrder),
            Include: true,
        }
    }
    if f.Step.Method != "POST" {
        return MidDecision{Class: ClassUpdate, Order: 900.0, Include: true}
    }
    if f.ParentResource == "" {
        return MidDecision{Class: ClassTopLevelCreate, Order: -1.0, Include: true}
    }
    if f.HasFirstTransition {
        return MidDecision{
            Class: ClassNestedUnderTransition,
            Order: f.FirstTransition + 0.5,
            Include: true,
        }
    }
    return MidDecision{Class: ClassNestedOrphan, Order: -0.5, Include: true}
}
```

**최대 depth = 2** (stateOps 안 branchSkip 체크). 기준점 이내.

### 2.3 파일 배치

```
pkg/generate/hurl/
├── decide_mid_step_class.go         신설 — StepClass/MidDecision/StepFacts + DecideMidStepClass + newStepFacts
├── decide_mid_step_class_test.go    신설 — 6 Class 전수 테이블 테스트 + 현행 로직 레퍼런스 비교
├── classify_mid_step.go             삭제 (호출처 1곳만 교체 후)
└── build_scenario_order.go          수정 — classifyMidStep 호출 교체
```

### 2.4 호출자 소비 패턴

```go
// build_scenario_order.go L72 근처
for _, s := range midCandidates {
    facts := newStepFacts(s, stateOps, branchSkip, transitionOrder, resourceFirstTransition)
    dec := DecideMidStepClass(facts)
    if !dec.Include {
        continue
    }
    midSteps = append(midSteps, orderedStep{step: s, order: dec.Order})
}
```

`classifyMidStep` 의 외부 호출처는 `build_scenario_order.go:72` 1곳. 교체 후 `classify_mid_step.go` 삭제. (교체 전 `grep -r classifyMidStep` 로 재확인.)

### 2.5 Toulmin 미적용 사유 (재확인)

현행이 이미 2-depth 이내 early-return chain. 기준점 이내 → if-else 유지.

**Toulmin 승격 조건** (장래):
- state 다이어그램 **사이클 지원** 도입 시 "return-to-state" warrant 가 추가되어 defeat 의미 발생
- 사용자 설정으로 classifier 우선순위 재정의가 필요할 때

둘 다 Phase010 범위 밖.

---

## 3. 검증

- `go test ./pkg/generate/hurl/...` — 6 Class 전수 테이블 테스트
- **레퍼런스 비교**: `decide_mid_step_class_test.go` 내에 현행 `classifyMidStep` 동등 함수를 복제하고, 임의 `StepFacts` 조합 1000개에 대해 `require.Equal` — bit-level 동일성 보증
- `fullend gen dummys/gigbridge/specs /tmp/x` — 생성된 `tests/hurl/*.hurl` step 순서가 기존과 bit-level 동일
- (가능시) `cd /tmp/x && hurl --test tests/hurl/*.hurl` 통과

---

## 4. 보류

- `sortDeletesByFK` / `topoSortDelete` / `buildTransitionOrder` / `authOrder` / `splitPrereqSteps` — 전부 **Toulmin 미적용**. 순수 그래프·정렬 알고리즘. Phase010 범위 밖.
- state 다이어그램 **사이클 지원** — 현행 BFS 가 visited 재방문 차단. 장래 확장 시 warrant 추가 고려 — Phase010 범위 밖.
- `scenarioStep` / `orderedStep` 타입 재설계는 Phase010 범위 밖.
