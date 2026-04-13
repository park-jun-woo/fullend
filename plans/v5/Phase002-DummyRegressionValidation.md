# Phase002 — Dummy 회귀 검증

## 목표

Phase001 이식 결과가 **internal 수준 산출물을 바이트 단위로 재현**함을 자동화된 스크립트로 검증한다.
gigbridge, zenflow 두 dummy 프로젝트를 기준 케이스로 삼아, 이식 전(베이스라인)과 이식 후(Phase001 결과)의 산출물을 전수 비교한다.

## 검증 수준 (Defence in Depth)

3겹 방어:

1. **파일 diff** — 바이트 동등성. 회귀 여부의 1차 감지.
2. **정적 빌드** — 생성된 Go 코드가 컴파일 되는가.
3. **문법 검증** — 생성된 Hurl 시나리오가 파싱되는가.

실행 환경(DB, 서버)까지 띄우는 실측 테스트는 **본 Phase 범위 밖**. 이는 Phase003+ 에서 별도 스위트로 다룸.

---

## 현황

| 항목 | 상태 |
|------|------|
| `dummys/gigbridge/specs/` | SSOT 완성, validate ✓ |
| `dummys/zenflow/specs/` | SSOT 완성, validate ✓ |
| `dummys/*/artifacts/` | **없음** — 지금까지 gen 실행한 적 없음 |
| 외부 도구: `hurl`, `oapi-codegen`, `sqlc` | 설치 확인됨 (`which` 성공) |
| 회귀 스크립트 | **없음** — 본 Phase에서 새로 작성 |
| git worktree 루트 | `master` 브랜치, 워크트리 1개 (`worktree-agent-afca9c16`) |

---

## 선행 작업 (Phase001 착수 전)

**핵심 제약**: 베이스라인은 Phase001 이 **시작되기 전**에 캡처되어야 한다.
Phase001이 시작되면 `internal/gen` 이 사라지므로, 당시의 산출물을 재현할 방법이 없다.

따라서 본 Phase 의 일부는 **시간적으로 Phase001 앞에** 배치된다.

### Step 0. Baseline 캡처 (Phase001 착수 직전)

1. `master` 에서 워킹 트리 clean 확인 (`git status` → 변경 없음).
2. 외부 도구 버전 고정:
   ```bash
   oapi-codegen --version > scripts/regression/tool-versions.txt
   sqlc version >> scripts/regression/tool-versions.txt
   hurl --version >> scripts/regression/tool-versions.txt
   go version >> scripts/regression/tool-versions.txt
   ```
3. 두 dummy 모두에 대해 `fullend gen` 실행 (internal 기반):
   ```bash
   go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
   go run ./cmd/fullend gen dummys/zenflow/specs  dummys/zenflow/artifacts
   ```
4. 생성 결과를 **baseline 디렉토리** 로 복사:
   ```bash
   cp -r dummys/gigbridge/artifacts dummys/gigbridge/artifacts.baseline
   cp -r dummys/zenflow/artifacts  dummys/zenflow/artifacts.baseline
   ```
5. 커밋:
   ```
   feat(regression): Phase002 baseline snapshot — gigbridge + zenflow
   ```

베이스라인은 **git에 commit** 하여 불변 레퍼런스로 보존. `.gitignore` 에 `artifacts/` 만 넣고 `artifacts.baseline/` 은 추적되게 한다.

---

## 변경 파일 목록

### 신규 — 스크립트

- `scripts/regression/README.md` — 사용법
- `scripts/regression/capture-baseline.sh` — 베이스라인 캡처 자동화
- `scripts/regression/run-regression.sh` — 이식 후 회귀 실행
- `scripts/regression/diff-artifacts.sh` — 바이트 diff + 요약 출력
- `scripts/regression/verify-build.sh` — 생성된 backend `go build` 검증
- `scripts/regression/verify-hurl.sh` — 생성된 smoke.hurl 문법 검증
- `scripts/regression/tool-versions.txt` — 캡처 시점 도구 버전 기록

### 신규 — 베이스라인 산출물

- `dummys/gigbridge/artifacts.baseline/` — gigbridge baseline (git tracked)
- `dummys/zenflow/artifacts.baseline/` — zenflow baseline (git tracked)

### 수정 — gitignore 조정

- `.ffignore` — **수정 금지** (프로젝트 규약)
- `.gitignore` — `dummys/*/artifacts/` 는 ignore, `dummys/*/artifacts.baseline/` 은 tracked (필요 시 엔트리 추가)

---

## 스크립트 설계

### `scripts/regression/capture-baseline.sh`

Phase001 착수 전 1회만 실행.

의사 코드:
```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

# 1. clean 상태 확인
if [[ -n "$(git status --porcelain)" ]]; then
    echo "working tree dirty — commit or stash first"; exit 1
fi

# 2. 도구 버전 기록
{
  echo "captured at: $(date -Iseconds)"
  echo "commit: $(git rev-parse HEAD)"
  echo "oapi-codegen: $(oapi-codegen --version 2>&1 | head -1)"
  echo "sqlc: $(sqlc version 2>&1)"
  echo "hurl: $(hurl --version 2>&1 | head -1)"
  echo "go: $(go version)"
} > scripts/regression/tool-versions.txt

# 3. 각 dummy 에 대해 gen + baseline 복사
for dummy in gigbridge zenflow; do
    specs="dummys/$dummy/specs"
    artifacts="dummys/$dummy/artifacts"
    baseline="dummys/$dummy/artifacts.baseline"

    rm -rf "$artifacts" "$baseline"
    go run ./cmd/fullend gen "$specs" "$artifacts"
    cp -r "$artifacts" "$baseline"
done

echo "baseline captured — commit dummys/*/artifacts.baseline + scripts/regression/tool-versions.txt"
```

### `scripts/regression/run-regression.sh`

Phase001 완료 후 실행.

의사 코드:
```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"

fail=0

for dummy in gigbridge zenflow; do
    specs="dummys/$dummy/specs"
    artifacts="dummys/$dummy/artifacts"
    baseline="dummys/$dummy/artifacts.baseline"

    echo "=== $dummy ==="
    rm -rf "$artifacts"
    go run ./cmd/fullend gen "$specs" "$artifacts"

    # 1. byte-level diff
    if ! scripts/regression/diff-artifacts.sh "$baseline" "$artifacts"; then
        echo "[FAIL] $dummy: artifact diff"
        fail=1
    fi

    # 2. build generated backend
    if ! scripts/regression/verify-build.sh "$artifacts/backend"; then
        echo "[FAIL] $dummy: backend build"
        fail=1
    fi

    # 3. verify hurl syntax
    if ! scripts/regression/verify-hurl.sh "$artifacts/tests"; then
        echo "[FAIL] $dummy: hurl syntax"
        fail=1
    fi
done

exit "$fail"
```

### `scripts/regression/diff-artifacts.sh`

```bash
#!/usr/bin/env bash
# usage: diff-artifacts.sh <baseline-dir> <current-dir>
set -euo pipefail
baseline="$1"
current="$2"

# tree diff (names only)
if ! diff -rq "$baseline" "$current"; then
    exit 1
fi
```

`diff -rq` 는 name + existence 만 비교. 파일 내용 diff 는 다르면 "differ" 표시로 감지됨. 자세한 content diff 는 실패 시 수동 실행 (`diff -ru`).

### `scripts/regression/verify-build.sh`

```bash
#!/usr/bin/env bash
# usage: verify-build.sh <backend-dir>
set -euo pipefail
backend="$1"

if [[ ! -f "$backend/go.mod" ]]; then
    echo "no go.mod in $backend — skipping build check"
    exit 0
fi

cd "$backend"
go build ./...
```

### `scripts/regression/verify-hurl.sh`

```bash
#!/usr/bin/env bash
# usage: verify-hurl.sh <tests-dir>
set -euo pipefail
tests="$1"

for f in "$tests"/*.hurl; do
    [[ -e "$f" ]] || continue
    # Hurl 은 --dry-run 이 없으므로 syntax 체크는 --test --retry-max 0 + 호스트 unreachable 로 간접 검증
    # 대신 파싱만 확인하는 방법: hurl --parse-only (hurl 4.3+)
    hurl --version | grep -q "hurl 4" && hurl --help | grep -q "no-output" || true
    # 문법만: 아래는 로컬 파싱 (실행 안 함)
    hurl --no-color --test --retry 0 --max-time 1 --variable host=http://127.0.0.1:1 "$f" 2>&1 \
        | grep -qE "(HTTP|Assertion|Connection refused|Could not resolve)" || {
            echo "hurl parse failed: $f"; exit 1
        }
done
```

Hurl 문법 체크는 버전에 따라 방법이 다름. 4.x 에서는 `--parse-only` 플래그가 없을 수 있어, **서버 무응답 + 연결 실패 까지 진행되면 문법은 OK** 로 간주. Phase 시작 시점에 설치된 hurl 버전으로 최종 확정.

---

## 작업 순서

### 선행 단계 (Phase001 앞)

1. **Step 0-a**: 회귀 스크립트 작성 (`scripts/regression/*.sh`) + 실행 권한 부여.
2. **Step 0-b**: `capture-baseline.sh` 실행. gigbridge·zenflow baseline 생성.
3. **Step 0-c**: 도구 버전 `tool-versions.txt` 기록 확인.
4. **Step 0-d**: baseline + 스크립트 커밋.

이 시점까지 완료 후에만 Phase001 착수.

### 후속 단계 (Phase001 완료 후)

5. **Step 1**: Phase001 마지막 커밋 기준 `run-regression.sh` 실행.
6. **Step 2**: 결과 분석:
   - 모든 dummy diff 0 + build 통과 + hurl 문법 OK → **Phase002 완료**.
   - 일부 실패 → Step 3.
7. **Step 3**: 차이 원인 분석 및 Phase001 수정:
   - **파일 내용 차이** — Phase001 이식 과정의 불완전 리네임·import 누락 가능. 수정 후 재커밋.
   - **트리 구조 차이** — 누락 생성기. Phase001 이식 범위 보강.
   - **build 실패** — 생성 코드의 import 경로·타입 시그니처 문제. Phase001 수정.
   - **hurl 문법 실패** — hurl 생성기 이식 중 `\n`·escape 처리 누락. 수정.
8. **Step 4**: Step 3 수정 후 Step 1 재실행. diff 0 도달까지 반복.
9. **Step 5**: 최종 통과 시 Phase002 완료 커밋:
   ```
   feat(regression): Phase002 통과 — gigbridge + zenflow diff 0
   ```

---

## 허용 오차 정책 (Tolerance Policy)

**원칙**: diff 0 이 기본. 어떤 차이도 허용하지 않음.

다만 아래의 **외부 요인** 차이는 별도 처리:

| 차이 유형 | 허용 여부 | 처리 |
|---------|---------|------|
| oapi-codegen 산출물 바이트 차이 (동일 버전 내) | 불허 | Phase001 에서 재현 실패 의미 |
| oapi-codegen 산출물 차이 (도구 버전 변경) | 허용 불가 — 버전 고정 필요 | `tool-versions.txt` 재검토 |
| sqlc 산출물 차이 (동일 버전 내) | 불허 | 동일 |
| 생성 시각·타임스탬프 | **없어야 함** | 산출물에 타임스탬프 삽입되면 Phase001 에서 제거 |
| 해시·contract 지시어 | 불허 | 입력 동일하면 해시 동일해야 함 |
| map 순회 비결정성 | 불허 | `sortedKeys` 등 정렬 유틸이 이미 보장 |

타임스탬프가 코드에 들어가는 경우(`// Code generated at: …` 같은 주석) 발견 시 Phase001 스펙 위반 — 즉시 수정 대상.

---

## 의존성

- **Phase001 계획** — 본 Phase 의 Step 0 은 Phase001 착수 전 완료 필수.
- **외부 도구** — `hurl`, `oapi-codegen`, `sqlc`, `go` 설치 확인됨.
- **dummy SSOT** — gigbridge + zenflow validate 통과 상태 유지 (이미 완료).

---

## 검증 방법

### 스크립트 자체 검증

- Step 0 완료 후 `dummys/gigbridge/artifacts.baseline/backend/` 에 파일이 생성되었는지 확인.
- `dummys/gigbridge/artifacts.baseline/tests/smoke.hurl` 존재 확인.
- `bash -n scripts/regression/*.sh` (syntax check) 통과.

### 회귀 자체 검증

- Phase001 완료 후 `scripts/regression/run-regression.sh` 종료 코드 0.
- 출력에 `[FAIL]` 이 하나도 없어야 함.
- 생성된 `dummys/*/artifacts/backend/` 에서 `go build ./...` 통과.
- 생성된 `dummys/*/artifacts/tests/*.hurl` 이 hurl 로 파싱 가능.

### 사후 확인

- `git status` 결과에 `dummys/*/artifacts/` 만 untracked (baseline 은 tracked).
- `tool-versions.txt` 가 최신 버전 반영.

---

## 주의사항

### 베이스라인 무결성

- baseline 은 Phase001 전에 단 1회 캡처. **이후 절대 재생성 금지**.
- 유일한 예외: 외부 도구 버전 의도적 변경 시. 이 경우 baseline 재캡처 + `tool-versions.txt` 갱신 + 별도 커밋으로 기록.

### 도구 버전 드리프트

- CI 또는 개발자 간 도구 버전 차이가 있으면 diff 가 false positive.
- 대응:
  - `tool-versions.txt` 를 회귀 실행 전 재확인.
  - 버전 불일치 시 스크립트가 즉시 중단하도록 확장 가능 (본 Phase 범위 밖 — Phase003 옵션).

### gitignore 재검토

현재 `.gitignore` 에 `dummys/*/artifacts/` 패턴이 있다면 `.baseline` 은 매칭 안 되게 구체화 필요:
```
dummys/*/artifacts/
!dummys/*/artifacts.baseline/
```

확인 후 필요 시 보강.

### `.ffignore` 금기

프로젝트 규약: `.ffignore` 수정 금지. baseline 디렉토리가 filefunc 스캔 대상이 되면 `.ffignore` 에 추가가 필요할 수도 있으나, **이는 별도 이슈** 로 분리. 본 Phase 는 `.ffignore` 건드리지 않음.

### dummy SSOT 변경 금지

Phase002 기간 중 gigbridge·zenflow SSOT 를 바꾸지 않는다. 변경 시 baseline 무의미해짐. SSOT 변경이 필요하면:
1. 본 Phase 통과 후 별도 변경
2. 변경 후 baseline 재캡처 + 새 커밋

### Hurl 실행 환경

`verify-hurl.sh` 는 존재하지 않는 호스트(`127.0.0.1:1`)로 시도하여 연결 실패까지 진행되면 문법 OK 로 간주. 이는 근사치 검증. 정확한 문법 검증은 hurl 버전 업그레이드 후 `--parse-only` 도입으로 개선 (향후 Phase).

### 워크트리 주의

현재 리포에 `worktree-agent-afca9c16` 가 있음. 회귀는 **반드시 메인 워크트리(`master`)** 에서 실행. 워크트리에서 실행 시 상대 경로 뒤섞일 수 있음.

---

## 완료 조건 (Definition of Done)

- [ ] `scripts/regression/*.sh` 작성 및 실행권한 부여
- [ ] `scripts/regression/README.md` 작성 (사용법)
- [ ] Phase001 착수 전 `capture-baseline.sh` 성공 실행
- [ ] `dummys/gigbridge/artifacts.baseline/`, `dummys/zenflow/artifacts.baseline/` 생성 및 git 커밋
- [ ] `scripts/regression/tool-versions.txt` 커밋
- [ ] Phase001 완료 후 `run-regression.sh` 종료코드 0
- [ ] 모든 dummy 에서 diff 0
- [ ] 모든 dummy 에서 backend `go build ./...` 통과
- [ ] 모든 dummy 에서 hurl 문법 검증 통과
- [ ] 최종 커밋: `feat(regression): Phase002 통과 — gigbridge + zenflow diff 0`

---

## 다음 Phase 예고

- **Phase003** — 복잡 로직 3군데 Toulmin 리팩토링:
  1. `pkg/generate/backend/generate_method_from_iface.go` 의 7-case switch
  2. `pkg/generate/backend/generate_main.go` 의 초기화 블록 조합
  3. `pkg/generate/hurl/build_scenario_order.go` + `classify_mid_step.go` 의 우선순위 그래프

  Phase003 각 서브스텝에서도 본 Phase 의 `run-regression.sh` 를 **리팩토링 게이트** 로 재사용한다.
  Toulmin 리팩토링은 본질적으로 "행동 보존" 이어야 하므로, 매 커밋마다 회귀 diff 0 유지 필수.

- **Phase004+** — 실제 실행 검증 (DB 기동, hurl 실제 요청) — 본 Phase 가 다루지 않은 영역.
