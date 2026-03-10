✅ 완료

# Phase 002: gluegen SSaC v2 마이그레이션

## 목표

SSaC v2 파서 타입 변경에 맞춰 `internal/gluegen/` 코드 생성기를 업데이트한다.

## 배경

gluegen은 SSaC 파서 출력을 소비하여 Go 백엔드 코드를 생성하는 최대 규모 모듈이다 (4,000+ LOC). SSaC v2의 원라인 DSL 타입 변경이 codegen 로직 전반에 영향을 미친다.

### 영향받는 파일

| 파일 | LOC | 변경 수준 |
|------|-----|----------|
| `gluegen/gluegen.go` | ~4,150 | 대 — 시퀀스 타입별 codegen 전면 |
| `gluegen/domain.go` | ~200 | 중 — ServiceFunc 접근 |
| `gluegen/model_impl.go` | ~500 | 중 — 모델 인터페이스 파생 |

## 변경 항목

### A. gluegen.go — 핵심 코드 생성

시퀀스 타입별 codegen 분기:

| v1 타입 | v2 타입 | codegen 변경 |
|---------|---------|-------------|
| `"get"` | `"get"` | Params → Args, 인자 생성 로직 |
| `"post"` | `"post"` | Params → Args |
| `"put"` | `"put"` | Params → Args |
| `"delete"` | `"delete"` | Params → Args |
| `"authorize"` | `"auth"` | @id → Inputs 맵 |
| `"guard nil"` | `"empty"` | Target 필드 (동일) |
| `"guard exists"` | `"exists"` | Target 필드 (동일) |
| `"guard state"` | `"state"` | DiagramID + Inputs |
| `"call"` | `"call"` | @func → Model, Params → Args |
| `"response json"` | `"response"` | @var → Fields 맵 |

주요 변경 패턴:
- 모든 `seq.Params` 순회 → `seq.Args` 순회
- `p.Source`, `p.Name` → `arg.Source`, `arg.Field`
- `p.Literal` 처리 → `arg.Literal` (동일 필드명이면 무변경)
- `"response json"` 분기 → `"response"` + `seq.Fields` 맵 사용

### B. domain.go — 도메인 폴더 구조

- `ServiceFunc.Domain` 접근 — 변경 없음 예상
- 시퀀스 순회 시 타입 상수만 업데이트

### C. model_impl.go — 모델 인터페이스 생성

- SSaC 시퀀스에서 Model.Method, Args를 추출하여 Go interface 메서드 시그니처 생성
- `seq.Params` → `seq.Args` 전환
- 인자 타입 매핑 로직 업데이트

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/gluegen/gluegen.go` | 시퀀스 타입 상수 + Params→Args 전면 교체 |
| `internal/gluegen/domain.go` | 타입 상수 업데이트 |
| `internal/gluegen/model_impl.go` | Params→Args, 인터페이스 생성 로직 |

## 의존성

- Phase 001 완료 (crosscheck 선행 필수는 아니나, 동일 타입 변경이므로 병렬 가능)
- SSaC v2 파서 로컬 참조

## 검증 방법

```bash
# 1. 빌드
go build ./internal/gluegen/...

# 2. fullend 전체 빌드
go build ./cmd/fullend/

# 3. dummy-lesson gen (Phase 003 이후 검증)
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
```
