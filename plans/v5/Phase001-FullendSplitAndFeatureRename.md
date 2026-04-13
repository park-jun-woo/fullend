# Phase001 — pkg/fullend 분리 + Domain→Feature 리네임

## 목표

파서 계열의 구조 정돈. 두 작업을 한 Phase 에 묶는다.

1. **`pkg/parser/fullend/` → `pkg/fullend/` 로 분리** — Fullstack 타입, ParseAll 함수, SSOT 탐지 유틸이 독립 패키지로 이동. `pkg/parser/` 는 개별 파서 전용으로 일관화.
2. **`pkg/parser/ssac.ServiceFunc` 의 `Domain` → `Feature` 필드 리네임** — 어휘 통일. 웹 도메인 의미(`pkg/parser/manifest/deploy.go`)는 보존.

**기존 소비자(validate/crosscheck)의 행동 변화 없음** — import 경로 갱신과 필드명 리네임만 적용.

검증 기준:
- `go build ./pkg/... ./internal/... ./cmd/...` 통과
- `go vet` 통과
- 기존 `go test ./pkg/...` 전부 통과

---

## 전제 및 비전제

### 전제

- `pkg/parser/fullend/fullstack.go` 에 Fullstack 이 이미 존재.
- `pkg/parser/fullend/parse_all.go` 에 `ParseAll(root, detected, skip) *Fullstack` 존재 확인.
- 같은 디렉토리에 `detect_ssots.go`, `detected_ssot.go`, `ssot_kind.go`, `parse_tangl_*.go`, `is_fullend_root.go`, `find_fullend_pkg_root.go` 등 **SSOT 탐지·orchestrator 보조 파일 일체** 존재 → Phase001 Part A 이동 범위에 전부 포함.
- `pkg/parser/ssac.ServiceFunc.Domain` 필드가 존재, 실측 pkg 내 참조 0건 (파서 내부 + 테스트만).

### 비전제

- internal/ 은 **건드리지 않는다** (복사 방식 유지).
- `pkg/ground/` 강화는 **Phase002** 에서 다룬다 (이 Phase 범위 밖).
- pkg/generate/ 는 이 Phase 에서 만들지 않는다 (Phase004 에서 이식).
- validate/crosscheck 는 이 Phase 에서 수정하지 않는다 (Phase003 에서 Ground 신 필드로 마이그).

---

## Part A — `pkg/fullend/` 분리

### 현황

```
pkg/parser/
├── ddl/ funcspec/ hurl/ manifest/ openapi/ rego/ ssac/ statemachine/ stml/
└── fullend/
    ├── fullstack.go     ← Fullstack 타입 정의
    └── parse_all.go     ← ParseAll() 함수
```

`pkg/parser/fullend/` 만 유일하게 **어그리게이터** — 나머지 형제는 단일 포맷 파서.
Fullstack 의 실제 소비자는 비파서 계열(validate/crosscheck/ground/generate). 위치 재배치가 개념상 맞다.

### 최종 배치

```
pkg/
├── fullend/                               ← 디렉토리 전체 이동
│   ├── fullstack.go
│   ├── parse_all.go
│   ├── detect_ssots.go
│   ├── detected_ssot.go
│   ├── ssot_kind.go
│   ├── parse_tangl_dir.go
│   ├── parse_tangl_file.go
│   ├── find_fullend_pkg_root.go
│   └── is_fullend_root.go
└── parser/
    ├── ddl/ funcspec/ hurl/ manifest/ openapi/ rego/ ssac/ statemachine/ stml/
    └── (fullend/ 제거됨)
```

즉 `pkg/parser/fullend/` 는 **단순 파서 집합이 아니라 SSOT 탐지·파싱·컨테이너 전부를 담는 어그리게이터 디렉토리**. 통째로 이동.

소비자 import 는 `pkg/fullend` 로 단축:
```go
import "github.com/park-jun-woo/fullend/pkg/fullend"

fs := fullend.ParseAll(specsDir, detected, skip)
validate.Run(fs)
```

### 작업 순서 (Part A)

**Step A-1.** 디렉토리 통째로 이동 (모든 파일 포함)
```
git mv pkg/parser/fullend pkg/fullend
```

**Step A-2.** 패키지 선언 유지 (`package fullend` 그대로 — 디렉토리명 일치).

**Step A-3.** 영향받는 import 경로 갱신. 영향 파일:
- `pkg/ground/build.go`
- `pkg/ground/populate_*.go` (Fullstack 참조가 있는 것들)
- `pkg/validate/**` 중 fullend import
- `pkg/crosscheck/**` 중 fullend import
- `internal/orchestrator/parsed.go` (혹시 pkg 참조가 있다면)
- 기타 grep 으로 식별

```
grep -rln "pkg/parser/fullend" pkg/ internal/ cmd/
# 발견된 모든 파일에서
# "github.com/park-jun-woo/fullend/pkg/parser/fullend" → "github.com/park-jun-woo/fullend/pkg/fullend"
```

**Step A-4.** 빌드 확인
```
go build ./pkg/... ./internal/... ./cmd/...
go vet ./pkg/... ./internal/... ./cmd/...
go test ./pkg/...
```

**Step A-5.** 커밋
```
refactor: Fullstack 을 pkg/fullend 로 분리 — pkg/parser 는 개별 파서 전용
```

---

## Part B — `Domain → Feature` 리네임 (파서)

Part A 와 같은 커밋에 묶어 처리 가능 (두 작업 모두 파서 계열 클린업).

### 대상

- `pkg/parser/ssac/service_func.go` — `Domain string` → `Feature string` 필드 리네임 + 주석 갱신.
- `pkg/parser/ssac/parse_dir_entry.go` — 할당부 + 에러 메시지.
- `pkg/parser/ssac/test_parse_domain_folder_test.go` → `test_parse_feature_folder_test.go` (파일명 + `TestParseDomainFolder` → `TestParseFeatureFolder`).
- `pkg/fullend/fullstack.go` — 필드 참조가 있다면.
- `pkg/validate/ssac/`, `pkg/crosscheck/`, `pkg/ground/` 의 `.Domain` 참조 — **실측 0건**, 확인만.

### 영향 없는 것

- `pkg/parser/manifest/deploy.go:7` 의 `Deploy.Domain` (웹 도메인). 리네임 대상 아님.
- `internal/ssac/parser` — Phase 범위 밖. `Domain` 유지.

### 커밋 묶음

Part A + Part B 를 하나의 커밋으로:
```
refactor: Fullstack 분리 + Domain→Feature 파서 리네임
```

---

## 작업 순서 요약

```
Part A + Part B 를 하나의 커밋으로 처리:
  - git mv pkg/parser/fullend pkg/fullend
  - pkg/parser/ssac/service_func.go 의 Domain → Feature
  - pkg/parser/ssac/parse_dir_entry.go 할당부 갱신
  - test_parse_domain_folder_test.go → test_parse_feature_folder_test.go 파일명 + 함수명
  - 영향받는 import 경로 일괄 갱신 (pkg/ground, pkg/validate, pkg/crosscheck 등)
  - 빌드·vet·테스트 통과 확인
  - 커밋: "refactor: Fullstack 분리 + Domain→Feature 파서 리네임"

전체 커밋 수: 1 개.
```

Ground 강화는 Phase002 에서 별도 수행.

---

## 주의사항

### R1. import 갱신 범위

`pkg/parser/fullend` 참조 위치 예상:
- pkg/ground/* (다수 — populate_*, build.go)
- pkg/validate/*, pkg/crosscheck/* (일부)
- internal/orchestrator/parsed.go (현재는 internal/genapi 사용 중, pkg 경로 미사용 가능성)

grep 으로 전수 파악 후 치환:
```
grep -rln "pkg/parser/fullend" pkg/ internal/ cmd/
```

### R2. 디렉토리 이동 범위 확인

`pkg/parser/fullend/` 에는 파서 단건 외에도 SSOT 탐지 (`detect_ssots.go`, `ssot_kind.go`, `detected_ssot.go`), 특수 파싱 (`parse_tangl_*.go`), 루트 식별 (`is_fullend_root.go`, `find_fullend_pkg_root.go`) 등 있음. **누락 없이 통째로 이동**.

### R3. ffignore 미수정

프로젝트 규약 — `.ffignore` 건드리지 않음. 파일 이동으로 filefunc 규칙 위반 경보가 유발되면 규칙 쪽을 고치지 말고 이동한 파일 내부를 조정.

### R4. 테스트 fixture 경로

테스트가 `pkg/parser/fullend` 에서 로드하던 것이 `pkg/fullend` 로 바뀌므로 테스트 import 도 갱신.

---

## 의존성

- 없음 (내부 리팩토링).
- `pkg/parser/ssac`, `pkg/parser/stml`, `pkg/parser/manifest` 존재 (확인됨).

---

## 완료 조건 (Definition of Done)

- [ ] `pkg/fullend/` 디렉토리 존재 (fullstack.go, parse_all.go, detect_ssots.go 등 원본 파일 전부)
- [ ] `pkg/parser/fullend/` 디렉토리 제거됨
- [ ] 모든 `pkg/parser/fullend` 참조가 `pkg/fullend` 로 교체
- [ ] `pkg/parser/ssac.ServiceFunc.Domain` → `Feature` 리네임 완료
- [ ] 웹 도메인 의미(`pkg/parser/manifest/deploy.go`) 보존
- [ ] `pkg/validate/**`, `pkg/crosscheck/**` 는 **코드 변경 없이** 기존 테스트 통과
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] `go test ./pkg/...` 통과

---

## 다음 Phase

- **Phase002** — `pkg/ground/` 를 generator 요구에 맞춰 강화 (구조적 신 필드 + populate 추가).
- **Phase003** — validate/crosscheck 를 확장된 Ground 신 필드로 점진 마이그.
- **Phase004** — internal 코드젠을 pkg/generate 로 복사 이식 (generator 는 Ground 신 필드 소비).
- **Phase005** — dummy 회귀 검증 (baseline 캡처 + diff).
