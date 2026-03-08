# ✅ Phase 009: sqlc + oapi-codegen Go import 전환

## 목표

sqlc와 oapi-codegen을 외부 도구 exec 호출에서 Go 라이브러리 import로 전환한다.
설치 불필요, 버전 go.mod 고정, pre-check 제거.

## 도구 및 라이센스

| 도구 | import 경로 | 라이센스 | 방식 |
|---|---|---|---|
| sqlc | `github.com/sqlc-dev/sqlc/pkg/cli` | MIT | `cli.Run([]string{"generate", ...})` |
| oapi-codegen | `github.com/oapi-codegen/oapi-codegen/v2` | Apache-2.0 | 라이브러리 API |
| terraform | (변경 없음) | BSL 1.1 | 외부 도구 exec 유지 |

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `go.mod` | sqlc, oapi-codegen 의존성 추가 |
| `artifacts/internal/orchestrator/gen.go` | genSqlc: exec → cli.Run() 전환, genOpenAPI() 함수 신규, pre-check에서 sqlc 제거 |

## sqlc 전환

기존:
```go
res := RunExec("sqlc", "generate", "-f", filepath.Join(specsDir, "sqlc.yaml"))
```

변경:
```go
import sqlccli "github.com/sqlc-dev/sqlc/pkg/cli"

func genSqlc(specsDir string) reporter.StepResult {
    step := reporter.StepResult{Name: "sqlc"}
    // cli.Run은 exit code 반환, stdout/stderr는 os 기본 출력
    code := sqlccli.Run([]string{"generate", "-f", filepath.Join(specsDir, "sqlc.yaml")})
    if code != 0 {
        step.Status = reporter.Fail
        step.Errors = append(step.Errors, fmt.Sprintf("sqlc generate failed (exit %d)", code))
        return step
    }
    step.Status = reporter.Pass
    step.Summary = "DB models generated"
    return step
}
```

## oapi-codegen 도입

oapi-codegen 라이브러리 API로 types + server 생성:

```go
func genOpenAPI(specsDir, artifactsDir string) reporter.StepResult {
    step := reporter.StepResult{Name: "oapi-gen"}
    apiPath := filepath.Join(specsDir, "api", "openapi.yaml")
    outDir := filepath.Join(artifactsDir, "backend", "api")
    os.MkdirAll(outDir, 0755)

    // 라이브러리 API로 types, server 생성
    // 실제 API 시그니처는 구현 시 확인 후 조정
    ...

    step.Status = reporter.Pass
    step.Summary = "types + server generated"
    return step
}
```

생성 위치:
```
<artifacts-dir>/
└── backend/
    └── api/
        ├── types.gen.go      # Go struct (request/response)
        └── server.gen.go     # net/http 핸들러 인터페이스 + 라우터
```

## gen 파이프라인 순서

```
1. pre-check    ← terraform만 검사 (선택, 경고 후 스킵)
2. sqlc         ← DB 모델 + 쿼리 구현 (import)
3. oapi-gen     ← OpenAPI 타입 + 서버 스텁 (import) ★ NEW
4. ssac-gen     ← 서비스 함수
5. ssac-model   ← 모델 인터페이스
6. stml-gen     ← React TSX
7. terraform    ← HCL 포맷팅 (외부 도구, 선택)
```

## pre-check 변경

sqlc, oapi-codegen 제거 (import이므로 설치 검사 불필요).
terraform만 남음 → terraform도 경고 후 스킵이므로 pre-check 자체가 필수 차단 없음.

```go
// pre-check: terraform은 경고 후 스킵
terraformAvailable := true
if _, ok := has[KindTerraform]; ok {
    if _, err := exec.LookPath("terraform"); err != nil {
        terraformAvailable = false
    }
}
// sqlc, oapi-codegen은 import이므로 검사 불필요
```

## 의존성

- ssac, stml 변경 없음
- go.mod에 sqlc, oapi-codegen 추가

## 검증 방법

1. `go build ./artifacts/cmd/... ./artifacts/internal/...` 성공
2. `fullend gen specs/dummy-lesson/ artifacts/dummy-lesson/` 실행 시:
   - `✓ sqlc  DB models generated`
   - `✓ oapi-gen  types + server generated`
3. `artifacts/dummy-lesson/backend/api/types.gen.go`, `server.gen.go` 생성 확인
4. sqlc, oapi-codegen CLI 미설치 상태에서도 정상 동작 확인
