# Phase022: 코드젠 버그 4건 수정 (zenflow 벤치마크에서 발견)

## 목표

dummy-zenflow 벤치마크 중 발견된 코드젠/문서 이슈 4건을 수정하여, SSOT → gen → build 파이프라인의 정합성을 높인다.

## 발견 경위

`files/dummy-zenflow-report2.md` — zenflow 벤치마크 (13분 31초, 4회 반복 수정 필요)

## 이슈 목록

### 이슈 1: `@empty` on int64 필드 → `== nil` 코드 생성 (HIGH)

**현상**: SSaC `@empty cr.Balance "msg" 402`에서 `Balance`가 `int64`일 때, 코드젠이 `cr.Balance == nil` 비교를 생성하여 컴파일 에러 발생.

**원인**: `@empty` 가드의 nil 비교 코드젠이 타입을 구분하지 않고 항상 `== nil`을 생성.

**수정 방안**:
- `@empty` 코드젠 시 대상 필드의 Go 타입을 확인
- pointer/slice/map/interface → `== nil`
- int/int64/float 등 값 타입 → `== 0` (zero-value 비교)
- string → `== ""`
- 또는 validate 단계에서 `@empty`가 값 타입 필드에 사용되면 WARNING/ERROR 발생

**변경 파일**:
- `internal/ssac/codegen.go` (또는 해당 코드젠 로직)

---

### 이슈 2: `@call` func에 모델 타입 배열 전달 시 타입 불일치 (MEDIUM)

**현상**: SSaC `@call worker.ProcessAction({Actions: actions})`에서 `actions`가 `[]model.Action`이지만, func spec의 `ProcessActionRequest.Actions`는 `[]worker.ActionItem`으로 선언됨. 코드젠이 타입 변환 없이 그대로 대입하여 컴파일 에러.

**원인**: `@call`에서 DDL 모델 변수를 func Request 필드에 전달할 때, 타입 호환성 검증/변환이 없음.

**수정 방안** (택 1):
- A) validate 단계에서 `@call` 인자의 타입과 func Request 필드 타입이 호환되는지 검증 → 불일치 시 ERROR
- B) codegen에서 동일 필드 구조의 타입 간 자동 변환 코드 생성 (JSON marshal/unmarshal 또는 필드별 복사)
- C) func spec에서 DDL 모델 타입을 직접 참조할 수 있도록 import 체계 확장

**변경 파일**:
- `internal/crosscheck/` (validate 쪽) 또는 `internal/ssac/codegen.go`

---

### 이슈 3: SSaC에서 정수 리터럴 사용 불가 — 매뉴얼에 미명시 (LOW)

**현상**: `@post ... CreditsSpent: 1` 작성 시, `1`이 변수명으로 해석되어 "변수가 선언되지 않았습니다" 에러.

**원인**: SSaC args 파서가 `"literal"` (문자열 리터럴)만 지원하고, 숫자 리터럴은 미지원.

**수정 방안** (택 1):
- A) 숫자 리터럴 지원 추가 — 파서에서 `[0-9]+`를 리터럴로 인식
- B) 매뉴얼에 "숫자 리터럴 미지원, func 응답으로 우회" 명시

**변경 파일**:
- A) `internal/ssac/parser.go` + `artifacts/manual-for-ai.md`
- B) `artifacts/manual-for-ai.md`만

---

### 이슈 4: func import 경로 제약 — 매뉴얼에 미명시 (LOW)

**현상**: SSaC에서 `import "github.com/org/project/func/billing"`으로 작성 시, `func-gen` 단계에서 "import 경로가 internal/ 또는 pkg/ 하위여야 합니다" 에러.

**원인**: func-gen이 import 경로에서 `internal/` 또는 `pkg/`를 기대하지만, specs의 물리 경로는 `func/`에 있어서 혼동 발생.

**수정 방안**:
- 매뉴얼 `artifacts/manual-for-ai.md`에 명시: "SSaC import 경로는 `internal/<pkg>` 형식으로 작성. func 스펙은 `specs/<project>/func/<pkg>/`에 위치하지만, SSaC import와 생성 코드에서는 `internal/<pkg>`로 참조됨"
- AGENTS.md에도 동일 내용 추가

**변경 파일**:
- `artifacts/manual-for-ai.md`
- `artifacts/AGENTS.md`

## 우선순위

1. 이슈 1 (HIGH) — 컴파일 에러 직결, 즉시 수정
2. 이슈 2 (MEDIUM) — 타입 안전성, validate에서 사전 차단
3. 이슈 3 (LOW) — 매뉴얼 보완 또는 파서 확장
4. 이슈 4 (LOW) — 매뉴얼 보완

## 검증 방법

1. zenflow 스펙으로 `fullend validate` → `fullend gen` → `go build` → `hurl --test` 통과
2. 이슈 1 재현: `@empty intField "msg"` → 컴파일 성공 확인
3. 이슈 2 재현: `@call func({ModelSlice: modelVar})` → validate에서 에러 또는 codegen에서 변환 확인
