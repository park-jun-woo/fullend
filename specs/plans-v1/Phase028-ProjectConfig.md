✅ 완료

# Phase 028: fullend.yaml 프로젝트 설정 + model import 수정

## 목표

specs 디렉토리에 `fullend.yaml`을 도입하여 프로젝트 메타데이터(모듈 경로, 언어/프레임워크 스택)를 선언한다.
gen 시 이 파일에서 Go module path, frontend package name 등을 읽어 산출물에 반영한다.
동시에 gluegen의 `model` import 누락 버그를 수정한다.

## 배경

- 현재 생성되는 `go.mod`의 module 경로가 출력 디렉토리 이름 그대로 들어감 (e.g. `module dummy-lesson-out/backend`)
- `*model.CurrentUser` 사용 시 `model` 패키지 import가 추가되지 않는 버그
- 언어/프레임워크 정보가 코드에 하드코딩되어 있어 확장 불가

## fullend.yaml 포맷

Kubernetes 스타일의 선언형 YAML. specs 디렉토리 루트에 위치.

```yaml
# fullend.yaml
apiVersion: fullend/v1
kind: Project

metadata:
  name: dummy-lesson

backend:
  lang: go
  framework: gin
  module: github.com/geul-org/dummy-lesson
  middleware:
    - bearerAuth

frontend:
  lang: typescript
  framework: react
  bundler: vite
  name: dummy-lesson-web

deploy:
  image: ghcr.io/geul-org/dummy-lesson
  domain: dummy-lesson.geul.org
```

### 필드 설명

| 필드 | 용도 | gen에서 사용처 |
|------|------|----------------|
| `apiVersion` | 포맷 버전 | 파서 호환성 |
| `kind` | 리소스 종류 (현재 `Project`만) | 향후 확장 |
| `metadata.name` | 프로젝트 식별자 | 디렉토리명 기본값 |
| `backend.lang` | 백엔드 언어 | 템플릿 선택 |
| `backend.framework` | 백엔드 프레임워크 | 템플릿 선택 |
| `backend.module` | Go module path | go.mod, import 경로 |
| `backend.middleware` | 미들웨어 목록 | `pkg/middleware/` 연결 |
| `frontend.lang` | 프론트엔드 언어 | 템플릿 선택 |
| `frontend.framework` | 프론트엔드 프레임워크 | 템플릿 선택 |
| `frontend.bundler` | 번들러 | vite.config 등 |
| `frontend.name` | npm 패키지명 | package.json |
| `deploy.image` | 컨테이너 이미지 | Terraform/Dockerfile |
| `deploy.domain` | 서비스 도메인 | Terraform/인프라 |

## 변경 항목

### A. fullend.yaml 파서 (`artifacts/internal/projectconfig/`)

- `projectconfig.go`: YAML 파서 + 구조체 정의
- `fullend.yaml` 없으면 validate 시 ERROR (필수 SSOT)
- `apiVersion: fullend/v1`, `kind: Project` 필수

### B. validate 연동 (`artifacts/internal/orchestrator/`)

- `fullend.yaml` 존재 여부 + 필수 필드 검증
- status 출력에 프로젝트 정보 표시

### B-2. crosscheck: middleware ↔ OpenAPI security 정합성

fullend.yaml의 `backend.middleware`와 OpenAPI `securitySchemes`의 일치를 검증한다.

| 규칙 | 조건 | 레벨 |
|------|------|------|
| OpenAPI `securitySchemes`에 이름 있음 → middleware에 동일 이름 필요 | 불일치 시 | ERROR |
| middleware에 이름 있음 → OpenAPI `securitySchemes`에 동일 이름 필요 | 불일치 시 | ERROR |
| OpenAPI 엔드포인트에 `security` 참조하는 이름이 middleware에 없음 | | ERROR |

### C. gluegen 수정 (`artifacts/internal/gluegen/`)

#### 1. modulePath를 fullend.yaml에서 읽기
- 현재: 출력 디렉토리명으로 go.mod 생성
- 변경: `backend.module` 값을 module path로 사용

#### 2. model import 누락 수정
- 현재: `QueryOpts{}` 패턴에서만 model import 추가
- 변경: `*model.CurrentUser` 패턴에서도 model import 추가
- `gluegen.go`의 `transformSource` 함수에서 model import 로직을 분리하여 모든 model 사용 패턴을 커버

#### 3. frontend package.json에 name 반영
- `frontend.name` → `package.json`의 `"name"` 필드

### D. dummy 스펙에 fullend.yaml 추가

#### `specs/dummy-lesson/fullend.yaml`
```yaml
apiVersion: fullend/v1
kind: Project

metadata:
  name: dummy-lesson

backend:
  lang: go
  framework: gin
  module: github.com/geul-org/dummy-lesson
  middleware:
    - bearerAuth

frontend:
  lang: typescript
  framework: react
  bundler: vite
  name: dummy-lesson-web
```

#### `specs/dummy-study/fullend.yaml`
```yaml
apiVersion: fullend/v1
kind: Project

metadata:
  name: dummy-study

backend:
  lang: go
  framework: gin
  module: github.com/geul-org/dummy-study
  middleware:
    - bearerAuth

frontend:
  lang: typescript
  framework: react
  bundler: vite
  name: dummy-study-web
```

### E. manual-for-ai.md 업데이트

- 디렉토리 구조에 `fullend.yaml` 추가
- fullend.yaml 문법 섹션 추가

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `artifacts/internal/projectconfig/projectconfig.go` | **신규** — YAML 파서 + 구조체 |
| `artifacts/internal/orchestrator/orchestrator.go` | fullend.yaml 로드 + validate 연동 |
| `artifacts/internal/crosscheck/crosscheck.go` | middleware ↔ OpenAPI security 정합성 검증 추가 |
| `artifacts/internal/gluegen/gluegen.go` | modulePath를 config에서 받기 + model import 수정 |
| `artifacts/internal/gluegen/domain.go` | modulePath 전달 경로 변경 |
| `specs/dummy-lesson/fullend.yaml` | **신규** — 프로젝트 설정 |
| `specs/dummy-study/fullend.yaml` | **신규** — 프로젝트 설정 |
| `artifacts/manual-for-ai.md` | fullend.yaml 문법 추가 |
| `go.mod` | `gopkg.in/yaml.v3` 의존성 추가 |

## 의존성

- `gopkg.in/yaml.v3` — YAML 파싱
- Phase 027 완료 (gin 전환)

## 검증 방법

```bash
# 1. fullend 빌드
go build ./artifacts/cmd/fullend/

# 2. validate — fullend.yaml 검증 포함
fullend validate specs/dummy-lesson
fullend validate specs/dummy-study

# 3. gen + build — module path 정상 반영 확인
fullend gen specs/dummy-lesson /tmp/gen-lesson
cd /tmp/gen-lesson/backend
head -1 go.mod  # → module github.com/geul-org/dummy-lesson
go build ./...  # → model import 에러 없음

# 4. dummy-study 동일
fullend gen specs/dummy-study /tmp/gen-study
cd /tmp/gen-study/backend
go build ./...
```
