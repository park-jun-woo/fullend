# Phase006 — Flat mode 제거 + pkg/generate 활성화 + orchestrator 배선

## 목표

Phase005 의 Ground 수렴 위에서 **pkg/generate 를 실제 동작시킨다**.

1. **Flat mode 제거** — pkg/generate/gogin 의 `*WithDomains` 중복 제거.
2. **pkg/generate.Generate stub 해소** — 실제 산출물 생성 배선.
3. **orchestrator 배선 교체** — internal/gen → pkg/generate.

성공 기준:
- `go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts` 실행 시 artifact 생성 (품질은 Tier 3 허용).
- orchestrator 가 pkg/generate 만 호출 (internal/gen import 0건).
- `go build` + `go vet` + `go test ./pkg/...` 통과.

---

## 전제

- **Phase005 완료** — 내부 함수가 Ground 기반.
- pkg/generate.Generate 는 여전히 stub 상태 (Phase004 에서 만든 것).
- internal/gen 이 여전히 살아있고 orchestrator 가 호출 중.

---

## 범위

### 포함

- pkg/generate/gogin 에서 Flat mode 잔재 제거
  - Flat 전용 파일 삭제 (`generate_main.go`, `generate_server_struct.go`, `transform_service_files.go` 등)
  - `*WithDomains` suffix 제거 (이제 유일한 버전이니 suffix 불필요)
  - `hasDomains`/`hasFeatures` 분기 제거 — 항상 Feature 경로
- pkg/generate.Generate (최상위), gogin.Generate, react.Generate, hurl.Generate stub 을 실제 로직으로 교체
- orchestrator 의 gen_glue.go / gen_ssac.go / gen_stml.go / parsed.go 가 pkg/generate 를 호출하도록 배선

### 포함하지 않음

- 구조 정리 (매개변수·결정 분산·템플릿) — Phase007
- Toulmin 포인트 도입 — Phase007
- Dummy 품질 검증 — Phase008

---

## 작업 순서

### Step 1. Flat mode 제거

pkg/generate/gogin/ 에서:

1. Flat 전용 파일 식별 (`*WithDomains` 가 존재하는 쌍의 base 버전):
   - `generate_main.go` (Flat) vs `generate_main_with_domains.go` (Feature)
   - `generate_server_struct.go` vs `generate_server_struct_with_domains.go`
   - `transform_service_files.go` vs `transform_service_files_with_domains.go`
   - `generate_auth_stub.go` (Flat) vs `generate_auth_stub_with_domains.go` (Feature)
   - `auth.go` (Flat 전용 `CurrentUser`)
   - `main_template.go` (Flat 전용 템플릿)

2. Flat 전용 파일 삭제.

3. `*WithDomains` 파일의 suffix 제거:
   - `generate_main_with_domains.go` → `generate_main.go`
   - `*WithDomains` 함수명도 동일하게 rename
   - 변수·주석 정돈

4. `generate.go` (pkg/generate/gogin) 의 `hasDomains/hasFeatures` 분기 제거 — 항상 Feature 경로 직진.

### Step 2. pkg/generate/gogin/generate.go 실제 로직 작성

Step 1 완료 후 gogin 이 단일 경로가 되었으므로, Fullstack + Ground 를 받아 내부 함수 호출하는 실제 Generate 메서드 작성.

참고: 현재 stub:
```go
func (g *GoGin) Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
    return fmt.Errorf("...")
}
```

교체 후 (개요):
```go
func (g *GoGin) Generate(fs *fullend.Fullstack, ground *rule.Ground, cfg *Config) error {
    intDir := filepath.Join(cfg.ArtifactsDir, "backend", "internal")
    if err := os.MkdirAll(intDir, 0755); err != nil { return err }
    // ... 순차적으로 transform, server, main, auth, middleware 호출 ...
}
```

### Step 3. pkg/generate/react, hurl stub 해소

동일 방식 — stub 에러 대신 실제 내부 함수 호출.

### Step 4. pkg/generate/generate.go 검증

최상위 `Generate` 가 Fullstack + Ground + Config + STMLGenOutput 받아 3개 하위 generator 호출. 이미 Phase004 에서 뼈대 있음.

### Step 5. orchestrator 배선 교체

- `internal/orchestrator/gen_glue.go`:
  - `internal/gen` import → `pkg/generate`
  - `gen.Generate(parsed, cfg, stmlOut)` → `generate.Generate(fs, ground, cfg, stmlOut)`
  - Fullstack + Ground 전달 방식으로 조정
- `internal/orchestrator/gen_ssac.go`:
  - `internal/ssac/generator` 호출 → `pkg/generate/gogin/ssac` 호출
- `internal/orchestrator/gen_stml.go`:
  - `internal/stml/generator` 호출 → `pkg/generate/react/stml` 호출
- `internal/orchestrator/parsed.go`:
  - `genapi.ParsedSSOTs` 반환 → `(*fullend.Fullstack, *rule.Ground, error)` 반환
  - 또는 중간 어댑터 (호출자 수정 최소화) 도 고려

### Step 6. 빌드 + 실행 검증

- `go build ./pkg/... ./internal/... ./cmd/...`
- `go vet`
- `go test ./pkg/...`
- `go run ./cmd/fullend gen dummys/gigbridge/specs /tmp/gigbridge-out` 실행 — 에러 없이 완료되는가?
- 생성된 artifact 가 존재 (품질 문제 있어도 OK)

---

## 주의사항

### R1. Step 1 (Flat 제거) 의 원자성

파일 삭제·rename·suffix 제거가 한꺼번에 일어나야 중간 빌드 깨지지 않음. **단일 커밋** 권장. 실패 시 `git reset --hard`.

### R2. orchestrator 배선 교체 시점

Step 1 + Step 2~4 가 먼저 완료돼야 orchestrator 가 호출할 대상이 준비됨. Step 5 는 마지막에.

### R3. internal/gen 은 그대로 유지

Phase006 완료 시 internal/gen 은 **호출되지 않는 dead code**. 삭제는 별도 Phase (문서엔 "Phase00N 별도") 에서.

### R4. 부분 실패 허용

생성된 artifact 의 `go build` 가 실패하거나 일부 파일이 누락돼도 OK. Phase008 에서 검증. 본 Phase 는 "pkg/generate 경로가 **일단 돌아간다**" 까지.

---

## 완료 조건 (Definition of Done)

- [ ] pkg/generate/gogin 에서 Flat 전용 파일 삭제
- [ ] `*WithDomains` suffix 제거 완료
- [ ] pkg/generate/{gogin,react,hurl,}.Generate stub 해소
- [ ] orchestrator 가 pkg/generate 만 호출 (internal/gen import 0)
- [ ] `go build` + `go vet` + `go test ./pkg/...` 통과
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 가능 (품질 관계없이 종료 코드 0)
- [ ] 커밋: `feat(generate): Flat 제거 + stub 해소 + orchestrator 배선 교체`

---

## 다음 Phase

- **Phase007** — 구조 정리 + Toulmin 포인트 3군데 도입.
- **Phase008** — Dummy 실용 검증 + 구조 건전성 지표.
