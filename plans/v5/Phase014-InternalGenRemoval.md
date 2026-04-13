# Phase014 — InternalGenRemoval

> `internal/gen/*`, `internal/ssac/generator`, `internal/stml/generator`, `internal/genapi`, `internal/contract` 일괄 삭제. v5 로드맵의 "참조 유지" 단계 종료.

## 목표

pkg/generate 가 internal 을 완전히 대체한 시점에 internal 을 일괄 제거해 코드베이스를 절반으로 축소. 예상 감소: **224 Go 파일 → 약 462 → 238 파일** (pkg 만 잔존).

---

## 전제

- **Phase012 완료** — orchestrator 의 internal/* 의존 0건
- **Phase013 완료** — 생성 품질이 pkg 경로에서 실용 가능
- `go build ./pkg/... ./cmd/...` 완전 통과
- `fullend gen` 이 **pkg 경로만으로** 모든 산출물 정상 생성
- 구조 지표 안정적 (Phase011 대비 개선 유지 또는 중립)

---

## 삭제 대상

### A. internal/gen/ (전체)

```
internal/gen/
├── gogin/          → 이미 pkg/generate/gogin 로 이식 (Phase004)
├── hurl/           → 이미 pkg/generate/hurl 로 이식
├── react/          → 이미 pkg/generate/react 로 이식
└── (공용 util)    → pkg/generate 로 이식 확인
```

### B. internal/ssac/generator, internal/stml/generator

Phase016 통합 이전 기준 ssac/stml 루트는 `internal/ssac/`, `internal/stml/`. 이 중 `generator/` 서브디렉만 제거. 파서(`internal/ssac/parser/`) 등은 **보존 대상 여부 재확인** 후 진행.

### C. internal/genapi (전체)

공용 타입 허브. pkg/fullend / pkg/parser/ssac / pkg/contract 로 흡수되었는지 확인 후 삭제.

### D. internal/contract

`pkg/contract` 가 대체. orchestrator 전환(Phase012) 후 안전하게 제거.

---

## 작업 순서

### Step 1. 의존성 전수 조사

```bash
rg "internal/(gen|ssac/generator|stml/generator|genapi|contract)" \
  --glob '!internal/(gen|ssac/generator|stml/generator|genapi|contract)/**' \
  --glob '*.go'
```

**0 hit 이어야 진행**. hit 있으면 Phase012/013 미완결 — 복귀.

### Step 2. 영향 범위 백업 (safety)

```bash
git checkout -b phase014-backup-before-removal
git push origin phase014-backup-before-removal
git checkout master
```

롤백 가능성 확보.

### Step 3. 삭제 (diff 크다)

```bash
rm -rf internal/gen
rm -rf internal/ssac/generator
rm -rf internal/stml/generator
rm -rf internal/genapi
rm -rf internal/contract
```

### Step 4. 빌드 확인

```bash
go build ./pkg/... ./internal/... ./cmd/...
```

실패 시: 놓친 의존 존재 — 되살리거나 pkg 로 추가 이식.

### Step 5. 테스트 + gen 회귀

```bash
go vet ./...
go test ./pkg/...
go run ./cmd/fullend validate dummys/gigbridge/specs
go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
cd dummys/gigbridge/artifacts/backend && go build ./... && cd -

go run ./cmd/fullend gen dummys/zenflow/specs dummys/zenflow/artifacts
cd dummys/zenflow/artifacts/backend && go build ./... && cd -
```

Phase013 빌드 통과 유지 확인.

### Step 6. 지표 재측정

```bash
go run scripts/structural_metrics.go > reports/metrics-phase014.md
```

삭제 후 pkg 내부 파일 수·평균 매개변수·WithDomains=0·Decide* ≥ 3 유지 확인.

### Step 7. `filefunc validate` 재실행

baseline 37 위반 중 삭제로 해소된 건 확인. 남은 위반은 pkg 잔여. 이후 Phase015 (TemplateAndResiduals) 이월.

### Step 8. `internal/gen/*/README.md` 9개 처리

Phase001 이전 분석 문서. 삭제 or `docs/archive/` 이관.

---

## 주의사항

### R1. 삭제는 한 번에

파일 그룹별 단계 삭제 금지. 위의 5 디렉토리는 **같은 커밋**에서 제거. 중간 상태는 빌드 깨짐 위험.

### R2. 복구 경로 확보

Step 2 의 백업 브랜치 + 태그를 반드시 먼저. `git reset --hard` 없이 `git revert` 로 되돌릴 수 있어야 함.

### R3. `internal/ssac/`, `internal/stml/` 루트 보존

ssac/stml 루트의 **parser** / **validator** 등 서브디렉는 별도 판단. 이 Phase 는 `generator/` 만 대상. 루트 자체 제거는 다른 Phase 에서.

### R4. dummy 정상 동작 유지

삭제 전후 gigbridge/zenflow 빌드가 동일하게 통과해야 함. 회귀 발생 시 즉시 중단 + 복구.

### R5. 긴급 우회 금지

빌드 깨질 경우 "임시로 internal 파일 복구" 식 우회 금지. 원인 추적 후 pkg 쪽에 누락분 이식.

---

## 완료 조건 (Definition of Done)

- [ ] Step 1 조사에서 0 hit 확인
- [ ] 5 디렉토리 삭제 완료
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet ./...` 통과
- [ ] `go test ./pkg/...` 통과
- [ ] gigbridge/zenflow `fullend gen` + backend `go build` 양쪽 통과
- [ ] `reports/metrics-phase014.md` 생성 (파일 수 대폭 감소 확인)
- [ ] `internal/gen/*/README.md` 9개 처리 (삭제 또는 `docs/archive/`)
- [ ] 커밋: `feat(internal): internal/{gen,ssac/generator,stml/generator,genapi,contract} 일괄 삭제`

## 의존

- **Phase012 완료 필수** — orchestrator 의존 0 보장
- **Phase013 완료 권장** — 생성 품질 안정 후 진행이 안전

## 다음 Phase

- **Phase015** — TemplateAndResiduals (text/template 전환 + 잔여 internal 의존 정리 + filefunc 논의)
