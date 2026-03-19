# Phase 003: Preserve 모드 — Go AST 함수 단위 splice ✅ 완료

## 목표

`fullend gen` 기본 동작을 preserve 모드로 변경. `//fullend:preserve` 함수의 body를 보존하고, `--reset`으로 전체 초기화.

## 배경

Phase 002까지는 모든 함수가 `//fullend:gen`이므로 전체 덮어씀. 이 Phase에서 개발자가 `gen` → `preserve`로 변경한 함수의 body를 보존하는 핵심 로직을 구현한다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/contract/splice.go` | 신규 — Go AST 함수 단위 splice 엔진 |
| `internal/contract/splice_test.go` | 신규 — splice 테스트 |
| `internal/orchestrator/gen.go` | 수정 — gen 흐름에 preserve 로직 통합 |
| `cmd/fullend/main.go` | 수정 — `--reset` 플래그 + Y/n 확인 |

## 상세 설계

### splice 엔진

```go
package contract

type SpliceResult struct {
    Content  string    // splice된 최종 파일 내용
    Warnings []Warning // contract 변경 경고
}

type Warning struct {
    File        string
    Function    string
    OldContract string
    NewContract string
}

// Splice는 기존 파일과 새로 생성된 내용을 함수 단위로 병합한다.
func Splice(oldPath string, newContent string) (*SpliceResult, error)
```

### Splice 알고리즘

```
1. oldPath 파일을 Go AST로 파싱 → 함수별 (디렉티브, body) 추출
2. newContent를 Go AST로 파싱 → 함수별 (디렉티브, body) 추출
3. 각 함수에 대해:
   - old에 없음 (새 함수) → new 그대로
   - old가 gen → new로 교체 (덮어씀)
   - old가 preserve + contract 동일 → old body 보존
   - old가 preserve + contract 다름 → old body 보존 + Warning 추가
   - old에 있지만 new에 없음 (SSOT에서 삭제) → 제거
4. import 재계산: 최종 함수들이 사용하는 패키지 합집합
5. go/format으로 출력
```

### 파일 레벨 디렉티브 처리

파일 첫 줄 `//fullend:` 디렉티브가 있으면 파일 전체를 하나의 단위로:

- `gen` → 파일 전체 덮어씀
- `preserve` → 파일 전체 보존 (함수 레벨 디렉티브 없으면)

### gen 흐름 변경

```go
// orchestrator/gen.go
func Gen(specsDir, artifactsDir string, reset bool) error {
    if reset {
        // preserve 함수 개수 표시 + Y/n 확인
        count := countPreserveFuncs(artifactsDir)
        if count > 0 {
            fmt.Printf("⚠ --reset: preserve 함수 %d개가 초기화됩니다.\n", count)
            fmt.Print("계속하시겠습니까? (Y/n): ")
            // 사용자 입력 확인
        }
    }

    // 기존 gen 로직 실행 → newContent 생성
    // ...

    // preserve 모드: splice 적용
    if !reset {
        for _, file := range generatedFiles {
            if exists(file.Path) {
                result, _ := contract.Splice(file.Path, file.Content)
                file.Content = result.Content
                // warnings 출력
            }
        }
    }
}
```

### 충돌 파일 생성

contract 변경 시 `.new` 파일 생성:

```
artifacts/backend/internal/service/gig/create_gig.go      ← 기존 body 유지
artifacts/backend/internal/service/gig/create_gig.go.new   ← 새 계약 기준 생성 코드
```

개발자가 머지 후 `.new` 삭제, 원본의 `contract=` 해시 갱신.

## 의존성

- Phase 001 (`internal/contract` — Directive, Hash)
- Phase 002 (디렉티브가 부착된 생성 코드)

## 검증

```bash
go test ./internal/contract/...
go test ./...
```

1. **기본 gen**: preserve 함수 body 보존 확인
2. **--reset**: 모든 preserve → gen 전환, body 재생성 확인
3. **--reset Y/n**: n 입력 시 중단 확인
4. **충돌**: contract 변경 시 `.new` 파일 생성 + 경고 출력 확인
5. **새 함수**: SSOT에 추가된 함수는 정상 생성 확인
6. **삭제된 함수**: SSOT에서 제거된 함수는 artifacts에서도 제거 확인
7. **import 재계산**: preserve body가 사용하는 패키지가 import에 포함 확인
8. **혼재**: 한 파일에 gen + preserve 함수 공존 시 정상 동작
