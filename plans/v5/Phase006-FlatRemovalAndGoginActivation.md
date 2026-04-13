# Phase006 — Flat mode 제거 + gogin.Generate 활성화

## 목표

Phase005 Ground 수렴 위에서 **gogin 단독 활성화** 에 집중. react/hurl 및 orchestrator 배선은 별도 Phase 로 분리.

1. **Flat mode 제거** — pkg/generate/gogin 의 `*WithDomains` 중복 제거 및 Flat 전용 파일 삭제.
2. **gogin.Generate stub 해소** — Fullstack + Ground + Config 받아 실제 백엔드 코드 생성.
3. **policy 타입 처리** — 내부 함수가 기대하는 타입과 호환.

성공 기준:
- `go build ./pkg/generate/gogin/...` 통과.
- `go test ./pkg/generate/gogin/...` 기존 테스트 통과.
- pkg/generate/gogin.Generate 가 stub 에러 아닌 실제 로직.

**react/hurl 의 stub 해소는 Phase007**.
**orchestrator 배선 교체는 Phase008**.

---

## 전제

- **Phase005 완료** — 내부 함수가 Ground 기반.
- pkg/generate/gogin 빌드 통과 (stub 상태).

---

## 범위

### 포함

- Flat mode 제거:
  - Flat 전용 파일 삭제 (generate_main.go[Flat], generate_server_struct.go[Flat], transform_service_files.go[Flat], auth.go, main_template.go)
  - `*WithDomains` → base 이름 rename
  - `has_domains.go` 삭제
  - `WithDomains` 함수명 suffix 제거
- gogin.Generate stub → 실제 로직 (internal/gen/gogin/generate.go 기반, Flat 분기 제거 버전)

### 포함하지 않음

- react.Generate, hurl.Generate stub 해소 — Phase007
- pkg/generate.Generate (최상위) 실사용 활성화 — Phase007
- orchestrator 배선 교체 — Phase008
- 구조 정리·Toulmin — Phase009
- Dummy 검증 — Phase010

---

## 작업 순서

### Step 1. Flat 제거 (완료)

이미 수행됨. 다음 커밋 참조:
- `refactor(generate/gogin): Flat mode 제거`

### Step 2. policy 타입 결정

내부 함수(`generate_main.go` 등) 가 기대하는 타입:
- `[]*policy.Policy` (internal/policy)

fs.Policies 는 `[]*ast.Module`, fs.ParsedPolicies 는 `[]rego.Policy`.

옵션:
- **A. internal/policy 계속 사용** — generate.Generate 내부에서 ParsedPolicies 변환 또는 별도 어댑터
- **B. 내부 함수 시그니처를 rego.Policy 로 변경** — 더 깔끔하지만 변경 범위 큼

내 권장: **A** (최소 변경). 어댑터는 간단:

```go
import (
    internalpolicy "github.com/park-jun-woo/fullend/internal/policy"
    "github.com/park-jun-woo/fullend/pkg/parser/rego"
)

func toInternalPolicies(parsed []rego.Policy) []*internalpolicy.Policy {
    // 구조 복사
}
```

또는 fs.Policies (ast.Module) 에서 internal/policy 로 재파싱.

### Step 3. gogin.Generate 구현

현재 stub 을 internal 원본 로직으로 교체 (Flat 분기 제거):

```go
func (g *GoGin) Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
    intDir := internalDir(cfg.ArtifactsDir)
    if err := os.MkdirAll(intDir, 0755); err != nil { return err }

    // Config 추출
    var claims map[string]manifest.ClaimDef
    var secretEnv string
    var queueBackend, sessionBackend, cacheBackend string
    var fileConfig *manifest.FileBackend
    if fs.Manifest != nil {
        if fs.Manifest.Backend.Auth != nil {
            claims = fs.Manifest.Backend.Auth.Claims
            secretEnv = fs.Manifest.Backend.Auth.SecretEnv
        }
        if fs.Manifest.Queue != nil { queueBackend = fs.Manifest.Queue.Backend }
        if fs.Manifest.Session != nil { sessionBackend = fs.Manifest.Session.Backend }
        if fs.Manifest.Cache != nil { cacheBackend = fs.Manifest.Cache.Backend }
        fileConfig = fs.Manifest.File
    }

    if hasBearerScheme(fs.OpenAPIDoc) && len(claims) == 0 {
        return fmt.Errorf("OpenAPI has bearerAuth but no claims config")
    }

    models := collectModels(fs.ServiceFuncs)
    funcs := collectFuncs(fs.ServiceFuncs)
    policies := toInternalPolicies(fs.ParsedPolicies)

    // Feature mode (Flat 제거)
    if err := transformServiceFiles(intDir, fs.ServiceFuncs, models, funcs, cfg.ModulePath, fs.OpenAPIDoc); err != nil { return err }
    if err := attachServiceDirectives(intDir, fs.ServiceFuncs); err != nil { return err }
    if err := generateAuthStub(intDir, cfg.ModulePath, claims); err != nil { return err }
    if err := generateAuthIfNeeded(intDir, cfg.ModulePath, claims, secretEnv); err != nil { return err }
    if err := generateServerStruct(intDir, fs.ServiceFuncs, cfg.ModulePath, fs.OpenAPIDoc); err != nil { return err }
    if err := generateMain(cfg.ArtifactsDir, fs.ServiceFuncs, cfg.ModulePath, queueBackend, policies, sessionBackend, cacheBackend, fileConfig); err != nil { return err }

    // 공유
    modelIncludeSpecs := collectModelIncludes(fs.OpenAPIDoc, fs.ServiceFuncs)
    cursorSpecs := collectCursorSpecs(fs.OpenAPIDoc)
    if err := generateModelImpls(intDir, models, cfg.ModulePath, cfg.SpecsDir, fs.ServiceFuncs, modelIncludeSpecs, cursorSpecs); err != nil { return err }
    if err := attachTSXDirectives(cfg.ArtifactsDir); err != nil { return err }

    return nil
}
```

### Step 4. 빌드 + 테스트

- `go build ./pkg/generate/gogin/...`
- `go vet`
- `go test ./pkg/generate/gogin/...` — 기존 테스트 유지

---

## 주의사항

### R1. orchestrator 는 아직 internal 경로

이 Phase 완료 후에도 orchestrator/gen_glue.go 는 internal/gen.Generate 를 계속 호출. 실사용자 CLI 는 pkg/generate 를 쓰지 않음. Phase008 에서 교체.

### R2. 테스트 신규 작성 최소

gogin.Generate 단위 테스트는 작성하지 않음. orchestrator 레벨 통합은 Phase010 dummy 검증에서.

### R3. 부분 실패 허용

일부 내부 함수 type 이슈로 빌드 안 되면 주석·TODO 처리. 목표는 "기본 경로 작동", 완벽 아님.

---

## 완료 조건 (Definition of Done)

- [x] Flat 제거 완료 (Step 1)
- [ ] gogin.Generate 실제 로직 구현
- [ ] `go build ./pkg/generate/gogin/...` 통과
- [ ] `go vet` 통과
- [ ] 기존 `go test ./pkg/generate/gogin/...` 유지
- [ ] 커밋: `feat(generate/gogin): Generate 활성화 — Fullstack + Ground 기반`

---

## 다음 Phase

- **Phase007** — react.Generate + hurl.Generate stub 해소.
- **Phase008** — orchestrator 배선 교체.
- **Phase009** — 구조 정리 + Toulmin.
- **Phase010** — Dummy 실용 검증.
