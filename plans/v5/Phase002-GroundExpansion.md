# Phase002 — pkg/ground 강화 (generator 요구 기준)

## 목표

`pkg/ground/` 에 **구조적 신 필드** (`Models`, `Tables`, `Ops`, `DTOs` 등) 와 대응 populate 를 추가한다. Phase004 에서 이식될 generator 가 이 필드를 직접 소비하도록 설계.

**기존 필드 무변경** (`Lookup`, `Types`, `Pairs`, `Schemas`, `Config`, `Vars`, `Flags`) — validate/crosscheck 의 행동 보존.

---

## 전제

- **Phase001 완료** — `pkg/fullend/` 분리됨, `Feature` 리네임됨.
- `pkg/ground/build.go`, `populate_*.go` 30+ 존재 (validate/crosscheck 용).

## 비전제

- internal/ 수정 안 함.
- validate/crosscheck 소비 지점 변경 없음 (Phase003 에서 수행).
- generator 이식 없음 (Phase004).

---

## Step 1. 설계 선행 — SymbolTable 접근 패턴 전수 조사

Phase002 의 **가장 중요한 체크포인트**. 구현 전 반드시 수행.

### 조사 범위

`internal/ssac/generator/` 63개 파일 (테스트 포함) 에서 `SymbolTable` 을 쓰는 모든 패턴 추출:

```bash
grep -rn "st\.\|symbolTable\.\|SymbolTable{" internal/ssac/generator/ > /tmp/symtab-usage.txt
```

분류 축:
- 어떤 필드를 읽는가 (Models, DDLTables, Operations, Funcs, DTOs, RequestSchemas)
- 접근 깊이 (1단, 2단, 3단 중첩)
- 반복 vs 단건 조회
- 쓰는 맥락 (HTTP 생성 / Subscribe 생성 / 모델 인터페이스 / Handler struct)

### 산출물

**설계 문서 `plans/v5/Phase002-GroundDesign.md` 작성** — 다음 포함:

1. SymbolTable 접근 패턴 전수 목록 (함수별)
2. 각 패턴을 Ground 신 필드로 매핑한 표
3. 신 필드 Go 타입 정의안
4. populate 함수 목록과 각각의 입력 (Fullstack 필드)
5. 매핑 실패 사례 (있으면) 와 대응 방안

### 산출물 제출 → 리뷰 → 착수 승인 절차

본 문서 작성 후 사용자 리뷰 후 구현 착수. 설계 흠 발견 시 재작성.

---

## Step 2. 신 타입 정의

설계 문서 승인 후.

### 파일 배치

- `pkg/rule/ground.go` — `Ground` 구조체에 신 필드 추가 (기존 필드 밑에 append)
- `pkg/rule/ground_types.go` (신설) — `ModelInfo`, `TableInfo`, `OperationInfo`, `ColumnInfo`, `MethodInfo`, `FieldInfo`, `ParamInfo`, `SchemaInfo`, `FKInfo` 등

### Ground 구조 예시

```go
// pkg/rule/ground.go
type Ground struct {
    // 기존 (보존)
    Lookup  map[string]StringSet
    Types   map[string]string
    Pairs   map[string]StringSet
    Config  map[string]bool
    Vars    StringSet
    Flags   StringSet
    Schemas map[string][]string

    // 신규 — generator 요구 (Phase002-GroundDesign.md 참조)
    Models map[string]ModelInfo
    Tables map[string]TableInfo
    Ops    map[string]OperationInfo
    DTOs   StringSet
}
```

정확한 필드 구성은 설계 문서에서 확정.

### 커밋

```
feat(ground): Ground 에 구조적 신 필드 타입 정의
```

빌드 통과 (기존 소비 코드 무영향) 후 커밋.

---

## Step 3. populate 함수 구현

설계 문서의 populate 매핑대로 각 함수 구현.

### 예상 신 populate 파일

- `pkg/ground/populate_models.go` — `g.Models` 채움 (Fullstack 의 ServiceFuncs, DDLTables 조합)
- `pkg/ground/populate_tables.go` — `g.Tables` 채움 (Fullstack 의 DDLTables, DDLResults)
- `pkg/ground/populate_ops.go` — `g.Ops` 채움 (Fullstack 의 OpenAPIDoc)
- `pkg/ground/populate_dtos.go` — `g.DTOs` 채움 (ModelDir 의 Go struct, `@dto` 태그)

정확한 목록은 설계 문서에서 확정.

### 커밋 단위

**populate 함수 하나당 커밋 하나** 권장:
```
feat(ground): populate_models 추가
feat(ground): populate_tables 추가
feat(ground): populate_ops 추가
feat(ground): populate_dtos 추가
```

각 커밋마다:
- `go build ./pkg/...` 통과
- 단위 테스트 작성 (dummy 프로젝트 fixture 기반)
- 테스트 통과 확인

---

## Step 4. Build 통합

`pkg/ground/build.go` 에 신 populate 호출 추가:

```go
func Build(fs *fullend.Fullstack) *rule.Ground {
    g := &rule.Ground{ ... }
    // 기존 populate 유지
    populateOpenAPI(g, fs)
    ...
    // 신규 populate 추가
    populateModels(g, fs)
    populateTables(g, fs)
    populateOps(g, fs)
    populateDTOs(g, fs)
    return g
}
```

### 커밋

```
feat(ground): Build 에 신 populate 통합
```

---

## Step 5. 최종 검증

### 정적
- `go build ./pkg/... ./internal/... ./cmd/...` 통과.
- `go vet` 통과.
- `go test ./pkg/...` — **기존 테스트 전부 그대로 통과** (backward compat 보증).

### 신 필드 값 검증
- dummy 프로젝트 gigbridge 기준 Ground 를 빌드.
- `g.Models`, `g.Tables`, `g.Ops`, `g.DTOs` 가 기대값 보유.
- 필요 시 표준 단위 테스트 추가 (`pkg/ground/test_populate_models_test.go` 등).

---

## 주의사항

### R1. 설계 선행 필수

Step 1 을 건너뛰고 구현 먼저 시작하면 Phase004 에서 "generator 가 Ground 로 못 바꾸는 패턴 발견" 으로 되돌아옴. 전수 목록화 → 매핑 가능 여부 확인 → 설계 → 구현 순 엄수.

### R2. 기존 populate 는 유지

Phase002 에서는 **신 필드 추가만**. 기존 평탄 populate (`populate_symbol_table.go` 등) 은 제거하지 않음. Phase003 에서 validate/crosscheck 마이그 후 유지/제거 결정.

### R3. 메모리 증가 용인

Ground 에 필드가 늘어나면서 메모리 사용량 증가. dummy 규모에서는 문제없음 (MB 단위). 대규모 프로젝트에서 문제 되면 lazy 채움 전환 고려 (향후).

### R4. 기존 테스트 fixture

Ground 를 빌드하는 기존 테스트에서 신 필드가 빈 값 허용되는지 확인. 필요 시 테스트 보조 생성자 (`ground.NewForTest()`) 도입 검토 — 하지만 가능하면 기존 테스트 무수정.

### R5. 순환 의존 방지

`pkg/rule/ground_types.go` 는 **pkg/parser/* 를 import 하지 않는다**. Ground 는 rule 패키지 내부의 순수 타입. Fullstack 에서 populate 시점에 값만 복사.

예:
```go
// O — populate 에서 Fullstack 을 읽어 Ground 채움
func populateModels(g *rule.Ground, fs *fullend.Fullstack) { ... }

// X — Ground 타입이 Fullstack 필드 참조
type ModelInfo struct {
    Parser *ssac.ServiceFunc  // ← 이렇게 하면 pkg/rule 이 pkg/parser/ssac 의존
}
```

---

## 의존성

- Phase001 완료.
- `pkg/fullend`, `pkg/parser/ssac`, `pkg/parser/ddl`, `pkg/parser/openapi` 존재.

---

## 완료 조건 (Definition of Done)

- [ ] `plans/v5/Phase002-GroundDesign.md` 설계 문서 작성 + 리뷰 통과
- [ ] `pkg/rule/ground.go` 에 신 필드 추가됨
- [ ] `pkg/rule/ground_types.go` (또는 유사) 에 신 타입 정의됨
- [ ] `pkg/ground/populate_*.go` 신 populate 함수들 추가
- [ ] `pkg/ground/build.go` 에 신 populate 호출 통합
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] `go test ./pkg/...` 전부 통과 (기존 테스트 무영향)
- [ ] dummy 프로젝트 기준 신 필드가 기대값 보유

---

## 다음 Phase

- **Phase003** — validate/crosscheck 가 신 필드를 점진 소비 + 중복 평탄 populate 제거.
- **Phase004** — internal 코드젠 → pkg/generate 복사 이식. generator 는 Ground 신 필드 소비.
- **Phase005** — dummy 회귀 검증 (baseline 캡처 + diff).
