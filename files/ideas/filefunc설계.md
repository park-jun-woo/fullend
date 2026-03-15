# filefunc 설계

## 1. 개요

코드 구조를 강제하는 validate 도구. 권고가 아니라 룰 — 위반 시 ERROR, CI에서 팅긴다.

`gofmt`가 겉모습(포매팅) 독재라면, filefunc는 **구조** 독재. 설계를 바꿔야 통과한다.

fullend에 통합하여 `fullend validate code`로 실행한다.

---

## 2. 룰

### 2.1 함수 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| R1 | 파일 하나에 func 하나 (파일명 = 함수명) | ERROR | tree-sitter func 노드 count |
| R2 | nesting depth 2까지만 허용 (3 이상 금지) | ERROR | tree-sitter depth 측정 |
| R3 | func max 1000 lines | ERROR | line count |
| R4 | func 권고 100 lines | WARNING | line count |
| R5 | 어노테이션 필수 (`//fullend:func`) | ERROR | 어노테이션 유무 |
| R6 | 어노테이션 값은 코드북에 존재해야 함 | ERROR | codebook yaml 대조 |
| R7 | desc 필수 (`//fullend:desc`) | ERROR | 어노테이션 유무 |
| R8 | input/output 타입 명시 | ERROR | Go AST / tree-sitter |
| R9 | 정형 구조 강제 (CLI는 cobra 등) | ERROR | import 검사 |

### 2.2 파일 구조 룰

| # | 룰 | 위반 시 |
|---|---|---|
| F1 | type/struct/interface/const/var → `types.go`에 모아야 함 | ERROR |
| F2 | `init()` → `main.go`에만 허용 | ERROR |
| F3 | 메서드 → 1 file 1 method (`server_start.go`, `server_stop.go`) | ERROR |
| F4 | `_test.go` → 복수 func 허용 | 예외 |

### 2.3 검토 후보

- cyclomatic complexity 상한
- 파라미터 개수 상한 (N개 초과 시 struct로 묶어라)
- return 값 개수 상한

---

## 3. 어노테이션

### 3.1 형식

```go
//fullend:func feature=crosscheck type=rule source=SSaC target=OpenAPI
//fullend:desc SSaC 함수명↔OpenAPI operationId 양방향 정합성 검증
//fullend:calls check_response_fields, check_err_status
//fullend:uses CrossError, ServiceFunc
func CheckSSaCOpenAPI(funcs []ServiceFunc, st *SymbolTable, doc *openapi3.T, specs []FuncSpec) []CrossError {
```

### 3.2 어노테이션 종류

| 어노테이션 | 용도 | 필수 |
|---|---|---|
| `//fullend:func` | 메타데이터 (feature, type, source, target 등) | O |
| `//fullend:desc` | 자연어 설명 (1줄) | O |
| `//fullend:calls` | 호출하는 함수 목록 | 자동 생성 |
| `//fullend:uses` | 사용하는 타입 목록 | 자동 생성 |

### 3.3 정보 저장 원칙

| 정보 | 저장 위치 | 이유 |
|---|---|---|
| 메타 (feature, type 등) | 어노테이션 | 불변, 짧음 |
| desc | 어노테이션 | 불변, 1줄 |
| calls/uses | 어노테이션 | 정적 관계, 코드 변경 시만 업데이트 |
| history | whyso (별도) | 무한 증가, 컨텍스트 오염 방지 |

---

## 4. 코드북

프로젝트별 `codebook-*.yaml`로 허용 값을 정의. 미등록 값 사용 시 ERROR.

```yaml
feature: [crosscheck, validate, gen, parse, report, contract, orchestrate]
type: [rule, parser, validator, generator, handler, middleware, loader, util]
pattern: [rulebook, target-interface, symbol-table, error-collection]
level: [ERROR, WARNING, INFO]
ssot: [OpenAPI, DDL, SSaC, STML, States, Policy, Config, Scenario, Func, Model]
```

---

## 5. func 노드 그래프

### 5.1 노드 구조

```yaml
node:
  name: CheckSSaCOpenAPI
  input: [ServiceFuncs, SymbolTable, OpenAPIDoc, FuncSpecs]
  output: [CrossError[]]
  desc: "SSaC 함수명↔OpenAPI operationId 양방향 정합성 검증"
```

### 5.2 함수 분류

모든 func은 2종류뿐:

1. **command func** — args/flags 파싱 → domain func 호출 → 출력 (접착제)
2. **domain func** — 실제 로직

domain func 간 관계도 2가지뿐:
- **호출** (A가 B를 call)
- **데이터 흐름** (A의 output이 B의 input)

---

## 6. func chain

### 6.1 feature 단위 탐색

같은 feature 안에서만 input/output 타입 매칭으로 체인을 구성한다:

```
ParseAll() → ParsedSSOTs → Run() → []CrossError → Print()
```

```bash
fullend chain func CheckSSaCOpenAPI   # 이 함수의 데이터 흐름
fullend chain feature crosscheck      # crosscheck feature 전체 체인
```

feature = 줌 레벨. 다른 feature와의 접점은 경계 노드로만 표시, drill-down 가능.

기존 `go callgraph`와의 차이:
- callgraph: 모든 호출 정적 분석 → 수천 노드, 그래프 폭발
- func chain: 같은 feature 안, 계약 기반 연결 → 필요한 것만

### 6.2 whyso 연동

func = file이므로 함수 단위 변경 이력이 파일 단위로 정확히 떨어진다:

```bash
whyso history check_ssac_openapi.go = CheckSSaCOpenAPI 함수의 변경 이력
```

### 6.3 암묵적 커플링 검출

whyso history는 사용자 요청 단위로 함께 수정된 파일을 추적한다:

```bash
whyso coupling check_ssac_openapi.go

같은 요청에 함께 수정된 함수:
  check_response_fields.go  8회
  check_err_status.go       5회
  types.go                  4회
```

calls/uses에 명시적 관계가 없는데 coupling에서 반복 출현 → **숨은 의존성 신호**:
- 같은 비즈니스 규칙을 다른 각도에서 구현
- interface 없이 암묵적으로 format을 맞춤
- 버그가 항상 같이 터짐

자동 WARNING: "이 두 함수는 명시적 관계 없이 N회 함께 수정됨. 의존성을 명시하세요."

---

## 7. 컨텍스트 엔지니어링

### 7.1 3단계 탐색 모델

```
1. graph 탐색 (LLM 불필요, 알고리즘)  → 관련 func 노드 특정
2. 메타데이터 읽기 (LLM 불필요, grep) → name, input, output, desc로 판별
3. body 읽기 (LLM 필요)              → 수정 대상 func만 컨텍스트 투입
```

현재는 1~3 전부 LLM이 수행 → 컨텍스트 폭발.
filefunc는 1, 2를 알고리즘으로 처리 → LLM은 3단계에서만 → 컨텍스트 오염 구조적 차단.

### 7.2 학술 근거

- **"Lost in the Middle" (Stanford, 2024)** — 관련 정보가 중간에 있으면 성능 30%+ 하락
- **"Context Length Alone Hurts" (Amazon, 2025)** — 불필요한 토큰이 공백이어도 성능 13.9~85% 하락
- **"Context Rot" (Chroma Research)** — focused prompt > full prompt (모든 모델)

연구는 증명했지만, 코드를 구조적으로 쪼개서 필요한 것만 넣는 도구가 없었다. filefunc가 그 빈자리.

### 7.3 핵심 원리

사람이 코드를 읽는 속도는 느리지만 기억은 오래 간다.
AI는 빠르게 읽지만 컨텍스트가 유한하다.
**greppable metadata는 AI의 약점을 정확히 보완한다.**

---

## 8. LLM 자동 어노테이션

기존 코드를 안 건드리고 메타데이터만 씌운다:

1. `fullend annotate ./internal/` — LLM이 func 읽고 어노테이션 + desc 자동 생성
2. 사람이 리뷰/수정
3. `fullend validate code ./internal/` — 정합성 검증

비침투적 — 기존 코드 그대로, 메타데이터 레이어만 추가. 도입 저항 제로.

---

## 9. 레지스트리

### 9.1 유명 라이브러리 사전 정의

GitHub 10k+ 라이브러리의 public API 메타를 미리 작업:

```
fullend registry
├── gin-gonic/gin          ★ 79k
├── gorilla/mux            ★ 21k
├── jmoiron/sqlx           ★ 16k
```

`import` 시 func 메타 자동 편입 → 내 코드 ↔ 외부 라이브러리 입출력 계약 검증.

### 9.2 외부 vs 자체 코드

| 대상 | 전략 | filefunc 룰 강제 |
|---|---|---|
| 외부 라이브러리 | public API 메타만 추출 | X |
| 자체 코드 | filefunc 룰 전면 강제 | O |

### 9.3 생태계 접근 제어

`fullend validate code` 통과 못하면 레지스트리 등록 불가. 생태계 품질의 바닥을 올린다.

---

## 10. fullend 통합

```bash
fullend validate saas specs/gigbridge   # SSOT 정합성 (기존)
fullend validate code ./internal/       # 코드 구조 (filefunc)
fullend annotate ./internal/            # LLM 자동 어노테이션
fullend chain func <name>               # func 데이터 흐름 추적
fullend chain feature <name>            # feature 전체 체인
```

---

## 11. 미결 사항

- 어노테이션 접두사: `//fullend:func` 확정? 범용 이름?
- 언어별 주석 문법 대응 (Go `//`, Python `#`, JS `//`)
- 레지스트리 호스팅/배포 방식
- LLM 어노테이션 품질 보장 (사람 리뷰 필수? 자동 승인?)
