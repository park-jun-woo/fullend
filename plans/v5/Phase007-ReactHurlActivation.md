# ✅ Phase007 — react.Generate + hurl.Generate 활성화 (완료)

## 목표

Phase006 에서 gogin 을 활성화한 뒤, 나머지 generator 두 개 (react, hurl) 도 stub 해소.

1. **react.Generate stub 해소** — Fullstack + STMLGenOutput 받아 프론트엔드 글루 생성.
2. **hurl.Generate stub 해소** — Fullstack + Ground 받아 Hurl smoke 테스트 생성.
3. **pkg/generate.Generate (최상위) 검증** — 세 generator 호출 체인 연결.

성공 기준:
- `go build ./pkg/generate/...` 통과.
- `go test ./pkg/generate/...` 기존 테스트 유지.
- pkg/generate.Generate 가 stub 에러 없이 세 generator 를 호출.

**orchestrator 배선 교체는 Phase008**.

---

## 전제

- **Phase006 완료** — gogin.Generate 활성화.
- pkg/generate/react, pkg/generate/hurl 빌드 통과 (stub).
- pkg/generate/generate.go 는 세 generator 를 호출하는 뼈대 (이미 Phase004 에서 작성).

---

## 범위

### 포함

- react.Generate stub → 실제 (internal/gen/react/generate.go 기반)
- hurl.Generate stub → 실제 (internal/gen/hurl/generate.go 기반)
- Fullstack 의 StateDiagrams/Policies/ServiceFuncs 를 내부 함수에 전달하는 타입 어댑터

### 포함하지 않음

- orchestrator 배선 — Phase008
- 구조 정리·Toulmin — Phase009

---

## 작업 순서

### Step 1. react.Generate 활성화

현재 stub:
```go
func Generate(fs *fullend.Fullstack, cfg *Config, stmlOut *STMLGenOutput) error {
    return fmt.Errorf("...")
}
```

교체:
```go
func Generate(fs *fullend.Fullstack, cfg *Config, stmlOut *STMLGenOutput) error {
    var deps map[string]string
    var pages []string
    var pageOps map[string]string
    if stmlOut != nil {
        deps = stmlOut.Deps
        pages = stmlOut.Pages
        pageOps = stmlOut.PageOps
    }
    return generateFrontendSetup(cfg.ArtifactsDir, fs.OpenAPIDoc, deps, pages, pageOps)
}
```

단 `generateFrontendSetup` 내부가 기대하는 타입이 pkg 와 호환되는지 확인.

### Step 2. hurl.Generate 활성화

현재 stub:
```go
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
    return fmt.Errorf("...")
}
```

교체:
```go
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
    return generateHurlTests(fs.OpenAPIDoc, cfg.ArtifactsDir, cfg.SpecsDir,
        fs.StateDiagrams, fs.ServiceFuncs, fs.Policies)
}
```

`generateHurlTests` 의 policies 인자 타입 확인 필요:
- 현재 internal/policy.Policy 기대?
- pkg/parser/rego 또는 fs.Policies (ast.Module) 중 선택

Policy 타입 어댑터 또는 시그니처 변경.

### Step 3. pkg/generate.Generate 검증

이미 Phase004 에서 작성된 최상위 배선:
```go
func Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config, stmlOut *react.STMLGenOutput) error {
    // gogin → react → hurl 순차 호출
}
```

세 하위 generator 가 활성화됐으므로 실제 동작할 것. 테스트.

### Step 4. 빌드 + 테스트

- `go build ./pkg/generate/...`
- `go vet`
- `go test ./pkg/generate/...`

---

## 주의사항

### R1. Policy 타입 통일 전략

gogin 과 hurl 모두 policy 관련 타입을 쓰므로 전략 통일 필요:
- 둘 다 internal/policy.Policy 유지 (Phase006 에서 결정된 기조)
- 둘 다 pkg 기반 (더 큰 작업, 별도 Phase 로 분리)

Phase006 결정에 맞춰 진행.

### R2. 생성기가 orchestrator 호출 안 됨

Phase007 완료 후에도 사용자 명령 (fullend gen) 에선 pkg/generate 가 호출되지 않음. Phase008 이 배선 교체해야 실제 사용.

### R3. internal/* 참조 허용

함수 내부가 여전히 internal/contract, internal/funcspec, internal/policy 등을 import 해도 OK. 완전 분리는 장기 목표.

---

## 완료 조건 (Definition of Done)

- [ ] react.Generate 가 stub 에러 아닌 실제 로직
- [ ] hurl.Generate 가 stub 에러 아닌 실제 로직
- [ ] `go build ./pkg/generate/...` 통과
- [ ] `go vet` 통과
- [ ] `go test ./pkg/generate/...` 기존 테스트 유지
- [ ] 커밋: `feat(generate): react + hurl 활성화 — Fullstack 기반 실제 생성`

---

## 다음 Phase

- **Phase008** — orchestrator 배선 교체.
- **Phase009** — 구조 정리 + Toulmin 포인트.
- **Phase010** — Dummy 실용 검증.
