# Phase 002: 디렉티브 부착 — fullend gen 시 //fullend:gen 자동 삽입 ✅ 완료

## 목표

`fullend gen` 실행 시 생성되는 모든 Go/TSX 함수에 `//fullend:gen` 디렉티브를 자동 부착한다.

## 배경

Phase 001에서 만든 contract hash + directive 인프라를 사용하여, 코드젠 출력에 소유권 메타를 내장한다. 이 Phase까지는 기존 동작과 동일 (전체 덮어씀). preserve 인식은 Phase 003.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/gluegen/gluegen.go` | 수정 — Generate 흐름에 contract hash 계산 추가 |
| `internal/gluegen/service_gen.go` | 수정 — 핸들러 함수 생성 시 디렉티브 부착 |
| `internal/gluegen/model_impl.go` | 수정 — 모델 구현 함수 생성 시 디렉티브 부착 |
| `internal/gluegen/stategen.go` | 수정 — state machine 파일 생성 시 파일 레벨 디렉티브 부착 |
| `internal/gluegen/middlewaregen.go` | 수정 — middleware 함수 생성 시 디렉티브 부착 |
| `internal/gluegen/frontgen.go` | 수정 — TSX 컴포넌트 생성 시 디렉티브 부착 |

## 상세 설계

### 핸들러 함수 (SSaC → Go)

ssac gen이 생성하는 각 핸들러 `.go` 파일에 함수별 디렉티브 삽입:

```go
// gluegen이 ssac gen 결과물을 후처리
// 각 함수의 doc comment 위치에 //fullend:gen 삽입
directive := contract.Directive{
    Ownership: "gen",
    SSOT:      relativePath(specsDir, ssacFile),
    Contract:  contract.HashServiceFunc(sf),
}
```

### 모델 구현 (DDL → Go)

gluegen이 `model/<name>.go` 생성 시, 각 메서드에 디렉티브 부착:

```go
//fullend:gen ssot=db/gigs.sql contract=e1d9f2
func (m *gigModelImpl) Create(...) { ... }

//fullend:gen ssot=db/gigs.sql contract=b8c3a7
func (m *gigModelImpl) FindByID(...) { ... }
```

### State Machine (Mermaid → Go)

파일 레벨 디렉티브 (package 선언 위):

```go
//fullend:gen ssot=states/gig.md contract=f5b3a9
package gigstate
```

### TSX (STML → React)

```tsx
// fullend:gen ssot=frontend/gig_list.html contract=d4e5f6
export function GigListPage() { ... }
```

## 의존성

- Phase 001 (`internal/contract`)

## 검증

```bash
fullend gen specs/gigbridge/ artifacts/gigbridge/
```

1. 생성된 모든 `.go` 핸들러 파일에 `//fullend:gen` 디렉티브 존재
2. 생성된 모든 모델 구현 함수에 `//fullend:gen` 디렉티브 존재
3. state machine 파일에 파일 레벨 디렉티브 존재
4. contract hash가 SSOT 내용에 따라 결정적 (동일 입력 → 동일 해시)
5. `go test ./...` 통과
