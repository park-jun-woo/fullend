# Phase033: control 어노테이션 전수 부착

## 목표

filefunc 개정 룰(A9)에 따라 모든 `//ff:func` 어노테이션에 `control=` 값을 추가한다.

> `control=`은 filefunc가 A9 룰 레벨에서 직접 강제한다. codebook.yaml에는 등록하지 않는다 (A2/A8 경로와 별개).

## 배경

filefunc 개정 — Böhm-Jacopini(1966) 기반 `control=` 어노테이션 도입:
- `sequence` — 순차 실행. depth 1에 switch/loop 없음. Q3 상한 100줄
- `selection` — switch/select가 depth 1 주 제어 구조. Q3 상한 300줄
- `iteration` — for/range가 depth 1 주 제어 구조. Q3 상한 100줄

**1 func 1 control** — switch+loop 혼합 금지. 혼합 시 내부 제어를 별도 함수로 추출.

검증 룰:
- A9: 모든 func 파일에 `control=` 필수 (sequence/selection/iteration) → ERROR
- A10: `control=selection`이지만 depth 1에 switch 없음 → ERROR
- A11: `control=iteration`이지만 depth 1에 loop 없음 → ERROR
- A12: `control=sequence`이지만 depth 1에 switch/loop 있음 → ERROR
- A13: `control=selection`이지만 depth 1에 loop 있음 → ERROR
- A14: `control=iteration`이지만 depth 1에 switch 있음 → ERROR

> `control=` 대상은 `//ff:func` 파일만 해당. `//ff:type` 파일은 A9 대상 아님 (매뉴얼: "Func files must have `control=`"). method 파일은 `//ff:func`으로 어노테이션하므로 대상에 포함.

## 현황

`//ff:func` 파일 **247개**, `//ff:type` 파일 40개.

func 파일의 control 분류 (depth 1 기준):

| 패키지 | func 파일 | sequence | selection | iteration | mixed |
|---|---|---|---|---|---|
| ssac/validator | 47 | 15 | 3 | 29 | 0 |
| gen/gogin | 85 | 22 | 4 | 57 | 2 |
| gen/hurl | 51 | 16 | 3 | 32 | 0 |
| orchestrator | 62 | 23 | 0 | 39 | 0 |
| gen (루트) | 2 | 2 | 0 | 0 | 0 |
| **합계** | **247** | **78** | **10** | **157** | **2** |

> mixed 2건 — 어떤 control 값을 붙여도 A13 또는 A14 위반. Phase035에서 리팩토링 후 control 부착:
> - `generate_method_from_iface.go` — depth 1에 for(76행) + switch(90행). switch가 주 제어 → switch 앞 for를 별도 함수로 추출 후 `control=selection`
> - `generate_state_machine_source.go` — depth 1에 for 2개 + switch 1개. for가 주 제어 → switch를 별도 함수로 추출 후 `control=iteration`

## 설계

### 1단계: 어노테이션 부착

247개 `//ff:func` 파일 전체에 `control=` 값 추가.

**분류 기준:**
- depth 1에 switch/select가 있으면 → `control=selection`
- depth 1에 for/range가 있으면 → `control=iteration`
- 둘 다 없으면 → `control=sequence`
- 둘 다 있으면 (mixed) → Phase035에서 분리 선행 후 분류 (A13/A14가 혼합을 강제 차단)

**변경 예시:**
```go
// BEFORE
//ff:func feature=ssac-validate type=rule

// AFTER
//ff:func feature=ssac-validate type=rule control=sequence
```

### 2단계: A9~A14 검증

`filefunc validate` 실행 — control 값과 실제 코드 구조 불일치 시 ERROR.

## 변경 파일

- `//ff:func`이 있는 247개 파일에 `control=sequence/selection/iteration` 추가

## 검증

1. `go build ./...` — 어노테이션만 변경이므로 빌드에 영향 없음
2. `filefunc validate` — A9~A14 위반 0 (mixed 2건은 Phase035 선행 분리 후 통과)
3. `go test ./...` — 코드 변경 없으므로 통과

## 리스크

- **분류 오류** — AST 기반 자동 분류. filefunc A10~A14가 즉시 검출.
- **Phase035 선행 의존** — mixed 2건(`generate_method_from_iface.go`, `generate_state_machine_source.go`)을 Phase035에서 먼저 분리해야 Phase033 위반 0 달성 가능. Phase035 → Phase033 순서 또는 동시 진행 필요.
